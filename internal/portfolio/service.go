package portfolio

import (
	"context"
	"fmt"

	"github.com/mhdalfrzh/stocklab-backend/internal/yahoofinance"

	"github.com/shopspring/decimal"
)

// YahooFetcher is the dependency interface for fetching stock quotes.
// The real *yahoofinance.Client satisfies this interface; tests provide a mock.
type YahooFetcher interface {
	FetchQuotes(ctx context.Context, symbols []string) map[string]yahoofinance.FetchQuoteResult
}

// Service defines the contract for generating a portfolio
// dashboard for a given user.
type Service interface {
	// GetDashboard returns the full dashboard response for the user's portfolio.
	// Returns ErrPortfolioNotFound if the user has no portfolio.
	GetDashboard(userID int64) (*DashboardResponse, error)
}

type service struct {
	repo         Repository
	yahooFetcher YahooFetcher
}

// NewService creates a Service with the
// provided repository and the real Yahoo Finance client.
func NewService(
	repo Repository,
	yahooClient *yahoofinance.Client,
) Service {
	return &service{
		repo:         repo,
		yahooFetcher: yahooClient,
	}
}

// NewServiceWithFetcher creates a Service
// with an injectable YahooFetcher. Use this in tests to inject a mock.
func NewServiceWithFetcher(
	repo Repository,
	fetcher YahooFetcher,
) Service {
	return &service{
		repo:         repo,
		yahooFetcher: fetcher,
	}
}

// GetDashboard orchestrates fetching portfolio data, market quotes, and
// computing the dashboard metrics.
func (s *service) GetDashboard(userID int64) (*DashboardResponse, error) {
	// 1. Fetch the user's portfolio.
	portfolio, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 2. Fetch holdings with stock info (single JOIN query).
	holdings, err := s.repo.FindHoldingsWithStocksByPortfolioID(portfolio.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch holdings: %w", err)
	}

	// 3. Empty portfolio — return zeroed summary with cash balance.
	if len(holdings) == 0 {
		return emptyDashboardResponse(portfolio.CashBalance), nil
	}

	// 4. Collect unique Yahoo symbols and batch-fetch quotes concurrently.
	symbols := make([]string, 0, len(holdings))
	for _, h := range holdings {
		symbols = append(symbols, h.YahooSymbol)
	}

	quoteResults := s.yahooFetcher.FetchQuotes(context.Background(), symbols)

	// 5. Build per-holding responses and accumulate summary values.
	var (
		holdingResponses         []DashboardHoldingResponse
		warnings                 []string
		totalStockMarketValue    = decimal.Zero
		totalCostOfAllHoldings   = decimal.Zero
		totalGainLoss            = decimal.Zero
		totalPreviousMarketValue = decimal.Zero
		totalCurrentMarketValue  = decimal.Zero
		hundred                  = decimal.NewFromInt(100)
	)

	for _, h := range holdings {
		result, ok := quoteResults[h.YahooSymbol]
		if !ok || result.Err != nil {
			// Record the warning and skip this holding from calculations.
			warnings = append(warnings, h.Ticker)
			continue
		}

		qty := decimal.NewFromInt(h.Quantity)
		currentPrice := result.Data.CurrentPrice
		previousClose := result.Data.PreviousClose

		// Per-holding metrics.
		marketValue := qty.Mul(currentPrice)
		totalCost := qty.Mul(h.AverageBuyPrice)
		gainLoss := marketValue.Sub(totalCost)

		var gainLossPct decimal.Decimal
		if !totalCost.IsZero() {
			gainLossPct = gainLoss.Div(totalCost).Mul(hundred).Round(2)
		}

		previousMV := qty.Mul(previousClose)

		// Accumulate summary values.
		totalStockMarketValue = totalStockMarketValue.Add(marketValue)
		totalCostOfAllHoldings = totalCostOfAllHoldings.Add(totalCost)
		totalGainLoss = totalGainLoss.Add(gainLoss)
		totalPreviousMarketValue = totalPreviousMarketValue.Add(previousMV)
		totalCurrentMarketValue = totalCurrentMarketValue.Add(marketValue)

		holdingResponses = append(holdingResponses, DashboardHoldingResponse{
			Ticker:             h.Ticker,
			CompanyName:        h.CompanyName,
			Quantity:           h.Quantity,
			AverageBuyPrice:    h.AverageBuyPrice.Round(2),
			CurrentPrice:       currentPrice.Round(2),
			MarketValue:        marketValue.Round(2),
			GainLoss:           gainLoss.Round(2),
			GainLossPercentage: gainLossPct,
		})
	}

	// 6. Compute summary-level metrics.
	totalPortfolioValue := portfolio.CashBalance.Add(totalStockMarketValue)

	var totalGainLossPct decimal.Decimal
	if !totalCostOfAllHoldings.IsZero() {
		totalGainLossPct = totalGainLoss.Div(totalCostOfAllHoldings).Mul(hundred).Round(2)
	}

	dayChange := totalCurrentMarketValue.Sub(totalPreviousMarketValue)

	var dayChangePct decimal.Decimal
	if !totalPreviousMarketValue.IsZero() {
		dayChangePct = dayChange.Div(totalPreviousMarketValue).Mul(hundred).Round(2)
	}

	// Ensure holdings slice is never nil in the JSON output.
	if holdingResponses == nil {
		holdingResponses = []DashboardHoldingResponse{}
	}

	response := &DashboardResponse{
		Summary: DashboardSummaryResponse{
			TotalPortfolioValue:     totalPortfolioValue.Round(2),
			CashBalance:             portfolio.CashBalance.Round(2),
			TotalGainLoss:           totalGainLoss.Round(2),
			TotalGainLossPercentage: totalGainLossPct,
			DayChange:               dayChange.Round(2),
			DayChangePercentage:     dayChangePct,
		},
		Holdings: holdingResponses,
		Warnings: warnings,
	}

	return response, nil
}

// emptyDashboardResponse returns the response for a portfolio with no holdings.
func emptyDashboardResponse(cashBalance decimal.Decimal) *DashboardResponse {
	return &DashboardResponse{
		Summary: DashboardSummaryResponse{
			TotalPortfolioValue:     cashBalance.Round(2),
			CashBalance:             cashBalance.Round(2),
			TotalGainLoss:           decimal.Zero,
			TotalGainLossPercentage: decimal.Zero,
			DayChange:               decimal.Zero,
			DayChangePercentage:     decimal.Zero,
		},
		Holdings: []DashboardHoldingResponse{},
	}
}

// compile-time assertion: *service satisfies Service.
var _ Service = (*service)(nil)

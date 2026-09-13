package portfolio

import "github.com/shopspring/decimal"

// DashboardHoldingResponse is the per-holding object in the dashboard response.
type DashboardHoldingResponse struct {
	Ticker             string          `json:"ticker"`
	CompanyName        string          `json:"companyName"`
	Quantity           int64           `json:"quantity"`
	AverageBuyPrice    decimal.Decimal `json:"averageBuyPrice"`
	CurrentPrice       decimal.Decimal `json:"currentPrice"`
	MarketValue        decimal.Decimal `json:"marketValue"`
	GainLoss           decimal.Decimal `json:"gainLoss"`
	GainLossPercentage decimal.Decimal `json:"gainLossPercentage"`
}

// DashboardSummaryResponse is the summary section of the dashboard response.
type DashboardSummaryResponse struct {
	TotalPortfolioValue     decimal.Decimal `json:"totalPortfolioValue"`
	CashBalance             decimal.Decimal `json:"cashBalance"`
	TotalGainLoss           decimal.Decimal `json:"totalGainLoss"`
	TotalGainLossPercentage decimal.Decimal `json:"totalGainLossPercentage"`
	DayChange               decimal.Decimal `json:"dayChange"`
	DayChangePercentage     decimal.Decimal `json:"dayChangePercentage"`
}

// DashboardResponse is the top-level dashboard API response.
type DashboardResponse struct {
	Summary  DashboardSummaryResponse `json:"summary"`
	Holdings []DashboardHoldingResponse `json:"holdings"`
	// Warnings lists tickers whose market data could not be retrieved.
	// Omitted from the JSON output when empty.
	Warnings []string `json:"warnings,omitempty"`
}

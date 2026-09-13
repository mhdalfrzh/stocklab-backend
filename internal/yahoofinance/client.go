// Package yahoofinance provides a client for fetching real-time and historical
// stock market data from the Yahoo Finance v8 chart API.
package yahoofinance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// Sentinel errors returned by the Yahoo Finance client.
var (
	// ErrSymbolNotFound is returned when Yahoo Finance responds with a 404
	// or returns no result for the requested symbol.
	ErrSymbolNotFound = errors.New("yahoofinance: symbol not found")

	// ErrUnavailable is returned when Yahoo Finance is temporarily unreachable
	// (network error, 5xx response, etc.).
	ErrUnavailable = errors.New("yahoofinance: service unavailable")

	// ErrNoMarketData is returned when Yahoo Finance returns a result but the
	// price fields are missing or zero.
	ErrNoMarketData = errors.New("yahoofinance: no market data in response")
)

// QuoteData holds the price data fetched from Yahoo Finance for one symbol.
type QuoteData struct {
	// CurrentPrice is the latest regular-market price.
	CurrentPrice decimal.Decimal
	// PreviousClose is the previous trading day's closing price.
	PreviousClose decimal.Decimal
}

// Client is the Yahoo Finance HTTP client.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a Client with sensible defaults.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    "https://query1.finance.yahoo.com",
	}
}

// NewClientWithBaseURL creates a Client pointing at a custom base URL.
// Used in tests to point at an httptest.Server.
func NewClientWithBaseURL(baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		baseURL:    baseURL,
	}
}

// --- Yahoo Finance JSON response shapes ---

type chartResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				RegularMarketPrice float64 `json:"regularMarketPrice"`
				ChartPreviousClose float64 `json:"chartPreviousClose"`
			} `json:"meta"`
		} `json:"result"`
		Error *struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"error"`
	} `json:"chart"`
}

// FetchQuote fetches the current price and previous close for a single symbol.
func (c *Client) FetchQuote(ctx context.Context, symbol string) (QuoteData, error) {
	url := fmt.Sprintf("%s/v8/finance/chart/%s?interval=1d&range=5d", c.baseURL, symbol)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return QuoteData{}, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	// Yahoo Finance returns 401 without a User-Agent header.
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return QuoteData{}, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return QuoteData{}, fmt.Errorf("%w: failed to read body: %v", ErrUnavailable, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return QuoteData{}, fmt.Errorf("%w: %s", ErrSymbolNotFound, symbol)
	}
	if resp.StatusCode >= 500 {
		return QuoteData{}, fmt.Errorf("%w: HTTP %d", ErrUnavailable, resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return QuoteData{}, fmt.Errorf("%w: HTTP %d", ErrUnavailable, resp.StatusCode)
	}

	var parsed chartResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return QuoteData{}, fmt.Errorf("%w: failed to parse JSON: %v", ErrUnavailable, err)
	}

	if parsed.Chart.Error != nil {
		if parsed.Chart.Error.Code == "Not Found" {
			return QuoteData{}, fmt.Errorf("%w: %s", ErrSymbolNotFound, symbol)
		}
		return QuoteData{}, fmt.Errorf("%w: %s", ErrUnavailable, parsed.Chart.Error.Description)
	}

	if len(parsed.Chart.Result) == 0 {
		return QuoteData{}, fmt.Errorf("%w: %s", ErrSymbolNotFound, symbol)
	}

	meta := parsed.Chart.Result[0].Meta
	if meta.RegularMarketPrice == 0 {
		return QuoteData{}, fmt.Errorf("%w: regularMarketPrice is zero for %s", ErrNoMarketData, symbol)
	}

	return QuoteData{
		CurrentPrice:  decimal.NewFromFloat(meta.RegularMarketPrice),
		PreviousClose: decimal.NewFromFloat(meta.ChartPreviousClose),
	}, nil
}

// FetchQuoteResult holds the result (or error) of fetching one symbol.
type FetchQuoteResult struct {
	Symbol string
	Data   QuoteData
	Err    error
}

// FetchQuotes fetches quotes for multiple symbols concurrently and returns
// a map of symbol → FetchQuoteResult for all requested symbols.
func (c *Client) FetchQuotes(ctx context.Context, symbols []string) map[string]FetchQuoteResult {
	results := make(map[string]FetchQuoteResult, len(symbols))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, sym := range symbols {
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			data, err := c.FetchQuote(ctx, symbol)
			mu.Lock()
			results[symbol] = FetchQuoteResult{Symbol: symbol, Data: data, Err: err}
			mu.Unlock()
		}(sym)
	}

	wg.Wait()
	return results
}

// Package domain contains the core GORM entity models that map directly to
// database tables. It must not import any HTTP or presentation-layer packages.
package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// Portfolio maps to the portfolios table.
type Portfolio struct {
	ID          int64           `gorm:"primaryKey"`
	UserID      int64           `gorm:"column:user_id;not null"`
	Name        string          `gorm:"column:name;not null"`
	CashBalance decimal.Decimal `gorm:"column:cash_balance;not null"`
	CreatedAt   time.Time       `gorm:"column:created_at"`
	UpdatedAt   time.Time       `gorm:"column:updated_at"`
}

func (Portfolio) TableName() string { return "portfolios" }

// Stock maps to the stocks table.
type Stock struct {
	ID          int64     `gorm:"primaryKey"`
	Ticker      string    `gorm:"column:ticker;not null"`
	YahooSymbol string    `gorm:"column:yahoo_symbol;not null"`
	CompanyName string    `gorm:"column:company_name;not null"`
	Sector      string    `gorm:"column:sector"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (Stock) TableName() string { return "stocks" }

// PortfolioHolding maps to the portfolio_holdings table.
type PortfolioHolding struct {
	ID              int64           `gorm:"primaryKey"`
	PortfolioID     int64           `gorm:"column:portfolio_id;not null"`
	StockID         int64           `gorm:"column:stock_id;not null"`
	Quantity        int64           `gorm:"column:quantity;not null"`
	AverageBuyPrice decimal.Decimal `gorm:"column:average_buy_price;not null"`
	CreatedAt       time.Time       `gorm:"column:created_at"`
	UpdatedAt       time.Time       `gorm:"column:updated_at"`
}

func (PortfolioHolding) TableName() string { return "portfolio_holdings" }

// HoldingWithStock is the flat result returned by the JOIN query in the
// portfolio repository. It contains all columns needed to build the dashboard
// without issuing further queries.
type HoldingWithStock struct {
	// From portfolio_holdings
	HoldingID       int64           `gorm:"column:holding_id"`
	Quantity        int64           `gorm:"column:quantity"`
	AverageBuyPrice decimal.Decimal `gorm:"column:average_buy_price"`
	// From stocks
	Ticker      string `gorm:"column:ticker"`
	YahooSymbol string `gorm:"column:yahoo_symbol"`
	CompanyName string `gorm:"column:company_name"`
}

package portfolio

import (
	"errors"

	"github.com/mhdalfrzh/stocklab-backend/internal/domain"

	"gorm.io/gorm"
)

// Repository defines the data access contract for portfolios.
type Repository interface {
	// FindByUserID returns the portfolio belonging to the given user.
	// Returns ErrPortfolioNotFound if no portfolio exists.
	FindByUserID(userID int64) (*domain.Portfolio, error)

	// FindHoldingsWithStocksByPortfolioID returns all holdings for the given
	// portfolio, joined with their associated stock information.
	// Returns an empty slice (not an error) when the portfolio has no holdings.
	FindHoldingsWithStocksByPortfolioID(portfolioID int64) ([]domain.HoldingWithStock, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new Repository backed by the
// provided GORM database connection.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// FindByUserID queries the portfolios table for the first portfolio owned by
// the given user. Returns ErrPortfolioNotFound when none exists.
func (r *repository) FindByUserID(userID int64) (*domain.Portfolio, error) {
	var p domain.Portfolio
	result := r.db.Where("user_id = ?", userID).First(&p)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrPortfolioNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &p, nil
}

// FindHoldingsWithStocksByPortfolioID performs a single JOIN query to retrieve
// all holdings and their associated stock data for the given portfolio. This
// avoids N+1 queries.
func (r *repository) FindHoldingsWithStocksByPortfolioID(portfolioID int64) ([]domain.HoldingWithStock, error) {
	var holdings []domain.HoldingWithStock

	result := r.db.Raw(`
		SELECT
			ph.id        AS holding_id,
			ph.quantity,
			ph.average_buy_price,
			s.ticker,
			s.yahoo_symbol,
			s.company_name
		FROM portfolio_holdings ph
		JOIN stocks s ON s.id = ph.stock_id
		WHERE ph.portfolio_id = ?
	`, portfolioID).Scan(&holdings)

	if result.Error != nil {
		return nil, result.Error
	}

	return holdings, nil
}

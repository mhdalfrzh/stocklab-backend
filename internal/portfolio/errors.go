package portfolio

import "errors"

// ErrPortfolioNotFound is returned when the user has no portfolio.
var ErrPortfolioNotFound = errors.New("portfolio not found for user")

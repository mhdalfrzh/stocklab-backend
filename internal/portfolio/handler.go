package portfolio

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests related to portfolios.
type Handler struct {
	service Service
}

// NewHandler creates a Handler with the given dashboard service.
func NewHandler(svc Service) *Handler {
	return &Handler{service: svc}
}

// GetDashboard handles GET /api/portfolios/dashboard
func (h *Handler) GetDashboard(c *gin.Context) {
	userID, ok := resolveUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing or invalid X-User-ID header",
		})
		return
	}

	dashboard, err := h.service.GetDashboard(userID)
	if err != nil {
		if errors.Is(err, ErrPortfolioNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "no portfolio found for this user",
			})
			return
		}

		// Unexpected error — return 500 but don't leak internal details.
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve portfolio dashboard",
		})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

// resolveUserID extracts the user ID from the X-User-ID request header.
// Returns (id, true) on success and (0, false) when the header is absent or
// not a valid integer.
func resolveUserID(c *gin.Context) (int64, bool) {
	raw := c.GetHeader("X-User-ID")
	if raw == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

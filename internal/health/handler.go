package health

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"bookstore-api/internal/shared/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) Handler {
	return Handler{db: db}
}

// Check godoc
// @Summary Check API health
// @Description Returns the current API health status.
// @Tags health
// @Produce json
// @Success 200 {object} response.Envelope{data=Status}
// @Failure 503 {object} response.Error
// @Router /health [get]
func (h Handler) Check(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "database is unavailable")
		return
	}

	response.OK(c, http.StatusOK, Status{Status: "ok"})
}

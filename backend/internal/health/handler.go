package health

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Health(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	err := h.db.PingContext(ctx)

	status := "ok"
	dbStatus := "up"

	if err != nil {
		status = "degraded"
		dbStatus = "down"
	}

	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"services": gin.H{
			"database": dbStatus,
		},
		"timestamp": time.Now(),
	})
}
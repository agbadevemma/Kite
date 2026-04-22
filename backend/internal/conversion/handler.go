package conversion

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Quote(c *gin.Context) {

	userID := c.GetString("userID")

	var req struct {
		From   string `json:"from_currency" binding:"required"`
		To     string `json:"to_currency" binding:"required"`
		Amount int64  `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	quote, err := h.service.Quote(c.Request.Context(), userID, req.From, req.To, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quote)
}

func (h *Handler) Execute(c *gin.Context) {

	var req ExecuteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quote_id required"})
		return
	}

	err := h.service.Execute(c.Request.Context(), req.QuoteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

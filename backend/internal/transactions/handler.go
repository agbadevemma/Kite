package transaction

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) List(c *gin.Context) {

	userID := c.GetString("userID")

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	txs, total, err := h.service.List(c.Request.Context(), userID, page, limit)
	if err != nil {
		log.Printf("Error fetching transactions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transactions"})
		return
	}

	if txs == nil {
		txs = []Transaction{}
	}
	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"data":       txs,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
	})
}

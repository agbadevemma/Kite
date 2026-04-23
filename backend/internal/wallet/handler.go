package wallet

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

type BalanceResponse struct {
	Currency string  `json:"currency"`
	Amount  float64 `json:"amount"`
}

func (h *Handler) GetBalances(c *gin.Context) {

	userID := c.GetString("userID")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	balances, err := h.service.GetBalances(c.Request.Context(), userID)

	
	if err != nil {
		log.Printf("Error fetching balances for user %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch balances"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balances": balances,
	})
}

package payout

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

type CreatePayoutDTO struct {
	Currency       string `json:"currency" binding:"required"`
	Amount         int64  `json:"amount" binding:"required"`
	Destination    string `json:"destination" binding:"required"`
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {

	var req CreatePayoutDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.GetString("userID")

	err := h.service.Create(
		c.Request.Context(),
		userID,
		req,
	)

	log.Printf("Payout creation error: %v", err)	

	if err != nil {

		if err.Error() == "insufficient balance" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create payout",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "payout created successfully",
	})
}
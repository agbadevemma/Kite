package deposit

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

type createDTO struct {
	Currency       string `json:"currency" binding:"required"`
	Amount         float64 `json:"amount" binding:"required"`
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {

	var body createDTO

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	userID := c.GetString("userID")

	err := h.service.Create(c.Request.Context(), Request{
		UserID:         userID,
		Currency:       body.Currency,
		Amount:         body.Amount,
		IdempotencyKey: body.IdempotencyKey,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
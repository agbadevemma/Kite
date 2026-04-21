package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type signupDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Signup(c *gin.Context) {
	var req signupDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.service.Signup(c.Request.Context(), SignupRequest(req))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "An error occured"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req loginDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := h.service.Login(c.Request.Context(), LoginRequest(req))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "An error occured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

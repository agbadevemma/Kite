package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/himarnoel/kite/internal/auth"
	"github.com/himarnoel/kite/internal/health"
)

func SetupRoutes(r *gin.Engine, db *sql.DB) {
	// INIT AUTH MODULE
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	// HEALTH CHECK
	healthHandler := health.NewHandler(db)
	r.GET("/health", healthHandler.Health)


	api := r.Group("/api")

	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/signup", authHandler.Signup)
		authRoutes.POST("/login", authHandler.Login)
	}


}

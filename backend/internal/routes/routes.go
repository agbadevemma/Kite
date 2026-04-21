package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/himarnoel/kite/internal/auth"
	"github.com/himarnoel/kite/internal/health"
	"github.com/himarnoel/kite/internal/ledger"
	"github.com/himarnoel/kite/internal/middleware"
	transaction "github.com/himarnoel/kite/internal/transactions"
	"github.com/himarnoel/kite/internal/wallet"
)

func SetupRoutes(r *gin.Engine, db *sql.DB) {
	// INIT  MODULE
	authRepo := auth.NewRepository(db)
	ledgerRepo := ledger.NewRepository(db)

	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)
	walletService := wallet.NewService(ledgerRepo)
	walletHandler := wallet.NewHandler(walletService)

	txRepo := transaction.NewRepository(db)
	txService := transaction.NewService(txRepo)
	txHandler := transaction.NewHandler(txService)

	// HEALTH CHECK
	healthHandler := health.NewHandler(db)
	r.GET("/health", healthHandler.Health)

	api := r.Group("/api")

	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/signup", authHandler.Signup)
		authRoutes.POST("/login", authHandler.Login)
	}

	walletRoutes := api.Group("/wallet")
	walletRoutes.Use(middleware.AuthMiddleware())
	{
		walletRoutes.GET("/balances", walletHandler.GetBalances)
	}

	api.GET("/transactions", middleware.AuthMiddleware(), txHandler.List)

}

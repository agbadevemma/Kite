package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/himarnoel/kite/internal/auth"
	"github.com/himarnoel/kite/internal/conversion"
	"github.com/himarnoel/kite/internal/deposit"
	"github.com/himarnoel/kite/internal/health"
	"github.com/himarnoel/kite/internal/ledger"
	"github.com/himarnoel/kite/internal/middleware"
	transaction "github.com/himarnoel/kite/internal/transactions"
	"github.com/himarnoel/kite/internal/wallet"
)

func SetupRoutes(r *gin.Engine, db *sql.DB) {
	// INIT  MODULE
	ledgerRepo := ledger.NewRepository(db)

	// HEALTH CHECK
	healthHandler := health.NewHandler(db)
	r.GET("/health", healthHandler.Health)

	api := r.Group("/api")

	// AUTH
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/signup", authHandler.Signup)
		authRoutes.POST("/login", authHandler.Login)
	}

	walletService := wallet.NewService(ledgerRepo)
	walletHandler := wallet.NewHandler(walletService)
	walletRoutes := api.Group("/wallet")
	walletRoutes.Use(middleware.AuthMiddleware())
	{
		walletRoutes.GET("/balances", walletHandler.GetBalances)
	}

	// TRANSACTIONS
	txRepo := transaction.NewRepository(db)
	txService := transaction.NewService(txRepo)
	txHandler := transaction.NewHandler(txService)

	api.GET("/transactions", middleware.AuthMiddleware(), txHandler.List)

	// DEPOSIT
	depositRepo := deposit.NewRepository(db)
	depositService := deposit.NewService(db, depositRepo, ledgerRepo, txRepo)
	depositHandler := deposit.NewHandler(depositService)

	api.POST("/deposits", middleware.AuthMiddleware(), depositHandler.Create)

	// CONVERSION
	convRepo := conversion.NewRepository(db)
	convService := conversion.NewService(convRepo, ledgerRepo, db)
	convHandler := conversion.NewHandler(convService)
	conversionRoutes := api.Group("/conversions")
	conversionRoutes.Use(middleware.AuthMiddleware())
	{
		conversionRoutes.POST("/quote", convHandler.Quote)
	conversionRoutes.POST("/execute", convHandler.Execute)
}
}

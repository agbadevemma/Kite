package main

import (
	"log"

	"github.com/joho/godotenv"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/himarnoel/kite/internal/db"
	"github.com/himarnoel/kite/internal/middleware"
	"github.com/himarnoel/kite/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	database := db.Connect()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://kite-frontend-di8t.onrender.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}), middleware.RateLimitMiddleware())

	routes.SetupRoutes(r, database)

	r.Run(":8080")
}

package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"github.com/himarnoel/kite/internal/db"
	"github.com/himarnoel/kite/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	database := db.Connect()

	r := gin.Default()

	routes.SetupRoutes(r, database)

	r.Run(":8080")
}

package main

import (
	"gin-crud/internal/models"
	"gin-crud/internal/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	models.InitDB()
	defer models.CloseDB()

	// Set up the Gin router with default middleware
	router := gin.Default()

	// Set up routes
	routes.SetupRoutes(router)

	// Start the server
	log.Println("Server started at http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

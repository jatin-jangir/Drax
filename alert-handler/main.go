package main

import (
	"alert-handler/db"
	"alert-handler/handlers"
	"alert-handler/middleware"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to MongoDB
	mongoClient, err := db.Connect("mongodb://localhost:27017")
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongoClient.Disconnect(nil)

	// Initialize Gin router
	r := gin.Default()

	// Enable CORS globally
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize handlers
	alertHandler := handlers.NewAlertHandler(mongoClient)
	authHandler := handlers.NewAuthHandler(mongoClient)

	// Public routes
	r.POST("/api/login", authHandler.Login)

	// Auth-protected routes
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/alerts", alertHandler.HandleAlert)
		auth.GET("/alerts", alertHandler.GetAlerts)
	}

	// Start server
	log.Println("Server running on :8080")
	r.Run(":8080")
}

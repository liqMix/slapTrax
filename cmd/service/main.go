package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var store *Store

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize BadgerDB
	dbPath := os.Getenv("BADGER_PATH")
	if dbPath == "" {
		dbPath = filepath.Join("data", "badger")
	}

	// Ensure directory exists
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	var err error
	store, err = NewStore(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize BadgerDB: %v", err)
	}
	defer store.Close()

	// Initialize router
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	// Public routes
	v1 := r.Group("/slapi/v1")
	v1.Use(gin.Logger())
	{
		// Health
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// Authentication routes
		v1.POST("/register", register)
		v1.POST("/login", login)
		v1.POST("/refresh", refresh)

		// Public leaderboard access
		v1.GET("/scores/leaderboard", getLeaderboard)
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(authMiddleware())
	{
		// User routes
		protected.GET("/user", getUser)
		// protected.POST("/users", updateUser)
		protected.GET("/user/scores", getUserScores)

		// Score routes
		protected.POST("/scores", createScore)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8008"
	}

	log.Printf("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

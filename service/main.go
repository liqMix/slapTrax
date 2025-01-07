package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=password dbname=holiday_db port=5435 sslmode=disable"
	}

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&User{}, &Score{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize router
	r := gin.Default()

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
	v1 := r.Group("/api/v1")
	{
		// Health
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// Authentication routes
		v1.POST("/register", register)
		v1.POST("/login", login)
		v1.POST("/refresh", refresh)
		v1.GET("/auth/check-auto-login", checkAutoLogin)
		v1.POST("/auth/auto-login", autoLogin)

		// Public leaderboard access
		v1.GET("/scores/leaderboard", getLeaderboard)
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(authMiddleware())
	{
		// User routes
		protected.GET("/users/:id", getUser)
		protected.PUT("/users/:id", updateUser)
		protected.GET("/users/:id/scores", getUserScores)

		// Score routes
		protected.POST("/scores", createScore)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

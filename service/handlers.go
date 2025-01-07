package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func getClientIP(c *gin.Context) string {
	// Check for X-Forwarded-For header first
	forwardedFor := c.GetHeader("X-Forwarded-For")
	if forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Fallback to RemoteAddr
	return c.ClientIP()
}

func checkAutoLogin(c *gin.Context) {
	clientIP := getClientIP(c)

	var user User
	result := db.Where("last_ip = ?", clientIP).
		Order("last_login_at DESC").
		First(&user)

	if result.Error != nil {
		c.JSON(http.StatusOK, AutoLoginResponse{Found: false})
		return
	}

	response := AutoLoginResponse{
		Found:    true,
		Username: user.Username,
		LastSeen: user.LastLoginAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}

func register(c *gin.Context) {
	var input struct {
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required,min=6"`
		DisplayName string `json:"display_name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := User{
		Username: input.Username,
	}

	if err := user.SetPassword(input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if result := db.Create(&user); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "username": user.Username})
}

func login(c *gin.Context) {
	var creds LoginCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if result := db.Where("username = ?", creds.Username).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := user.CheckPassword(creds.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Update IP and login time on successful login
	clientIP := getClientIP(c)
	user.LastIP = clientIP
	user.LastLoginAt = time.Now()

	accessToken, err := GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	if err := user.SetRefreshToken(refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save refresh token"})
		return
	}

	db.Save(&user)

	c.JSON(http.StatusOK, TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func autoLogin(c *gin.Context) {
	clientIP := getClientIP(c)

	var user User
	result := db.Where("last_ip = ?", clientIP).
		Order("last_login_at DESC").
		First(&user)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No previous login found for this IP"})
		return
	}

	// Update login time
	user.LastLoginAt = time.Now()

	accessToken, err := GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	if err := user.SetRefreshToken(refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save refresh token"})
		return
	}

	db.Save(&user)

	c.JSON(http.StatusOK, TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func refresh(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if result := db.Where("refresh_token != ''").First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	if err := user.CheckRefreshToken(input.RefreshToken); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	accessToken, err := GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	if err := user.SetRefreshToken(refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save refresh token"})
		return
	}

	// Update last login time
	user.LastLoginAt = time.Now()
	db.Save(&user)

	c.JSON(http.StatusOK, TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c.GetHeader("Authorization"))
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		claims, err := ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	var user User

	if result := db.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("userID")

	// Ensure users can only update their own data
	if userID != id {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update other users"})
		return
	}

	var user User
	if result := db.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updateData User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := db.Model(&user).Updates(updateData); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func createScore(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var score Score
	if err := c.ShouldBindJSON(&score); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	score.UserID = userID
	score.PlayedAt = time.Now()

	if result := db.Create(&score); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, score)
}

func getUserScores(c *gin.Context) {
	id := c.Param("id")
	var scores []Score

	if result := db.Where("user_id = ?", id).Order("score desc").Find(&scores); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, scores)
}

func getLeaderboard(c *gin.Context) {
	song := c.Query("song")
	difficulty := c.Query("difficulty")
	limit := 100 // Default limit for leaderboard entries

	var scores []struct {
		Score
		Username string `json:"username"`
	}

	query := db.Table("scores").
		Select("scores.*, users.username").
		Joins("left join users on users.id = scores.user_id").
		Where("song = ? AND difficulty = ?", song, difficulty).
		Order("score desc").
		Limit(limit)

	if result := query.Find(&scores); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, scores)
}

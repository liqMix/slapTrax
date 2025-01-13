package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
)

func getClientIP(c *gin.Context) string {
	forwardedFor := c.GetHeader("X-Forwarded-For")
	if forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	return c.ClientIP()
}

func register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required,min=1,max=8"`
		Password string `json:"password" binding:"required,min=1"`
		// Settings external.Settings `json:"settings" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &User{
		Username: input.Username,
		// Settings: input.Settings,
	}

	if err := user.SetPassword(input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if err := store.CreateUser(user); err != nil {
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

	user, err := store.GetUserByUsername(creds.Username)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := user.CheckPassword(creds.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Update IP and login time
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

	if err := store.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

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

	// Find user with non-empty refresh token
	var users []User
	err := store.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte(userPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			var user User
			err := it.Item().Value(func(v []byte) error {
				return json.Unmarshal(v, &user)
			})
			if err != nil {
				continue
			}
			if user.RefreshToken != "" {
				users = append(users, user)
			}
		}
		return nil
	})

	if err != nil || len(users) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Find user with matching refresh token
	var validUser *User
	for _, user := range users {
		if err := user.CheckRefreshToken(input.RefreshToken); err == nil {
			validUser = &user
			break
		}
	}

	if validUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate new tokens
	accessToken, err := GenerateAccessToken(validUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	if err := validUser.SetRefreshToken(refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	if err := store.UpdateUser(validUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

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
	id := c.MustGet("userID").(uint)

	user, err := store.GetUserByID(id)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.Clean()
	c.JSON(http.StatusOK, user)
}

// func updateUser(c *gin.Context) {
// 	id := c.MustGet("userID").(uint)

// 	user, err := store.GetUserByID(id)
// 	if err != nil || user == nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	// var settings external.Settings
// 	// if err := c.ShouldBindJSON(&settings); err != nil {
// 	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	// 	return
// 	// }

// 	// user.Settings = settings
// 	user.UpdatedAt = time.Now()

// 	if err := store.UpdateUser(user); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	user.Clean()
// 	c.JSON(http.StatusOK, user)
// }

func createScore(c *gin.Context) {
	id := c.MustGet("userID").(uint)

	user, err := store.GetUserByID(id)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var score Score
	if err := c.ShouldBindJSON(&score); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user score for this song/difficulty
	currentScores, err := store.GetUserScores(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ratingIncrease := getRatingValue(&score)
	for _, currentScore := range currentScores {
		if currentScore.SongHash == score.SongHash && currentScore.Difficulty == score.Difficulty {
			value := getRatingValue(&currentScore)
			ratingIncrease -= value
			break
		}
	}

	if ratingIncrease < 0 {
		ratingIncrease = 0
	}
	score.UserID = id
	score.Username = user.Username

	if err := store.CreateScoreAndUpdateRating(&score, id, ratingIncrease); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, score)
}

func getUserScores(c *gin.Context) {
	id := c.MustGet("userID").(uint)
	logger.Debug("Getting scores for user %d", id)
	_, err := store.GetUserByID(id)
	if err != nil {
		logger.Error("Failed to get user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	scores, err := store.GetUserScores(id)
	if err != nil {
		logger.Error("Failed to get user scores: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, scores)
}

func getLeaderboard(c *gin.Context) {
	song := c.Query("song")
	difficultyStr := c.Query("difficulty")

	var difficulty int
	_, err := fmt.Sscanf(difficultyStr, "%d", &difficulty)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid difficulty"})
		return
	}

	scores, err := store.GetLeaderboard(song, difficulty)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, score := range scores {
		user, err := store.GetUserByID(score.UserID)
		if err != nil {
			continue
		}
		score.Rank = user.Rank
	}

	c.JSON(http.StatusOK, scores)
}

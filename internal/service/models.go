package main

import (
	"time"

	"github.com/liqmix/ebiten-holiday-2024/internal/external"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string            `gorm:"unique;not null" json:"username"`
	Password     string            `json:"-"` // Password hash, not exposed in JSON
	Settings     external.Settings `gorm:"type:jsonb" json:"settings"`
	RefreshToken string            `json:"-"` // Stored hashed refresh token
	LastIP       string            `json:"-"` // Store last successful login IP
	LastLoginAt  time.Time         `json:"last_login_at"`
}

type Score struct {
	gorm.Model
	UserID     uint      `json:"user_id"`
	Username   string    `json:"username"`
	Song       string    `json:"song" gorm:"not null"`
	Score      int64     `json:"score" gorm:"not null"`
	Accuracy   float64   `json:"accuracy"`
	MaxCombo   int       `json:"max_combo"`
	PlayedAt   time.Time `json:"played_at" gorm:"not null"`
	Difficulty string    `json:"difficulty" gorm:"not null"`
}

type LoginCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AutoLoginResponse struct {
	Found    bool   `json:"found"`
	Username string `json:"username,omitempty"`
	LastSeen string `json:"last_seen,omitempty"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) SetRefreshToken(token string) error {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.RefreshToken = string(hashedToken)
	return nil
}

func (u *User) CheckRefreshToken(token string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.RefreshToken), []byte(token))
}

func (u *User) UpdateSettings(updates map[string]interface{}) error {
	// Remove protected fields
	delete(updates, "username")
	delete(updates, "password")

	// Only allow updating permitted fields
	allowedFields := []string{"settings"}
	for field := range updates {
		allowed := false
		for _, af := range allowedFields {
			if field == af {
				allowed = true
				break
			}
		}
		if !allowed {
			delete(updates, field)
		}
	}

	// Apply remaining updates
	if len(updates) > 0 {
		return db.Model(u).Updates(updates).Error
	}
	return nil
}

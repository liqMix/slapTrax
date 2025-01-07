package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string    `gorm:"unique;not null" json:"username"`
	Password     string    `json:"-"` // Password hash, not exposed in JSON
	Settings     Settings  `gorm:"type:jsonb" json:"settings"`
	RefreshToken string    `json:"-"` // Stored hashed refresh token
	LastIP       string    `json:"-"` // Store last successful login IP
	LastLoginAt  time.Time `json:"last_login_at"`
}

type Settings struct {
	Version            string  `json:"version"`
	Locale             string  `json:"locale"`
	Fullscreen         bool    `json:"fullscreen"`
	ScreenWidth        int     `json:"screen_width"`
	ScreenHeight       int     `json:"screen_height"`
	RenderWidth        int     `json:"render_width"`
	RenderHeight       int     `json:"render_height"`
	FixedRenderScale   bool    `json:"fixed_render_scale"`
	BGMVolume          float64 `json:"bgm_volume"`
	SFXVolume          float64 `json:"sfx_volume"`
	SongVolume         float64 `json:"song_volume"`
	LaneSpeed          float64 `json:"lane_speed"`
	AudioOffset        int64   `json:"audio_offset"`
	InputOffset        int64   `json:"input_offset"`
	WaveringLane       bool    `json:"wavering_lane"`
	NoteColorTheme     string  `json:"note_color_theme"`
	CenterNoteColor    string  `json:"center_note_color"`
	CornerNoteColor    string  `json:"corner_note_color"`
	DisableHoldNotes   bool    `json:"disable_hold_notes"`
	DisableHitEffects  bool    `json:"disable_hit_effects"`
	DisableLaneEffects bool    `json:"disable_lane_effects"`
}

type Score struct {
	gorm.Model
	UserID     uint      `json:"user_id"`
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

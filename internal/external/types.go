package external

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// LoginState represents the current authentication state
type LoginState int

const (
	StateUninitialized LoginState = iota
	StateOffline
	StateOnline
	StateLoggingIn
)

// Settings represents user preferences and configuration
type Settings struct {
	Version      int       `json:"version"`
	LastModified time.Time `json:"last_modified"`
	LastSync     time.Time `json:"last_sync"`
	IsNewUser    bool      `json:"is_new_user"`

	// Game Settings
	Locale           string `json:"locale"`
	Fullscreen       bool   `json:"fullscreen"`
	ScreenWidth      int    `json:"screen_width"`
	ScreenHeight     int    `json:"screen_height"`
	RenderWidth      int    `json:"render_width"`
	RenderHeight     int    `json:"render_height"`
	FixedRenderScale bool   `json:"fixed_render_scale"`
	KeyConfig        int    `json:"key_config"`

	// Audio Settings
	BGMVolume   float64 `json:"bgm_volume"`
	SFXVolume   float64 `json:"sfx_volume"`
	SongVolume  float64 `json:"song_volume"`
	AudioOffset int64   `json:"audio_offset"`
	InputOffset int64   `json:"input_offset"`

	// Gameplay Settings
	NoteWidth          float32 `json:"note_width"`
	LaneSpeed          float64 `json:"lane_speed"`
	NoteColorTheme     string  `json:"note_color_theme"`
	CenterNoteColor    string  `json:"center_note_color"`
	CornerNoteColor    string  `json:"corner_note_color"`
	DisableHoldNotes   bool    `json:"disable_hold_notes"`
	DisableHitEffects  bool    `json:"disable_hit_effects"`
	DisableLaneEffects bool    `json:"disable_lane_effects"`
	EdgePlayArea       bool    `json:"fullscreen_play_area"`
	Use3DNotes         bool    `json:"use_3d_notes"`
}

func (s *Settings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan implements sql.Scanner interface
func (s *Settings) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte for Settings, got %T", value)
	}

	return json.Unmarshal(bytes, s)
}

// User represents the current user state
type User struct {
	Username string    `json:"username"`
	Rank     float64   `json:"rank"`
	Settings *Settings `json:"settings"`
}

// TokenPair represents authentication tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Session represents an authenticated user session
type Session struct {
	Username     string    `json:"username"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type Score struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UserID     uint      `json:"user_id"`
	Username   string    `json:"username"`
	SongHash   string    `json:"song_hash"`
	Score      int       `json:"score"`
	Rank       float64   `json:"rank"`
	Accuracy   float64   `json:"accuracy"`
	MaxCombo   int       `json:"max_combo"`
	PlayedAt   time.Time `json:"played_at"`
	Difficulty int       `json:"difficulty"`
}

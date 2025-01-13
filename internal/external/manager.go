package external

import (
	"fmt"
	"net/http"
	"time"

	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
)

// Manager handles all user state including auth and settings
type Manager struct {
	currentUser *User
	session     *Session
	loginState  LoginState
	storage     *Storage
	rememberMe  bool
	isConnected func() bool
}

// New creates a new state manager
func NewManager(storagePath string) *Manager {
	return &Manager{
		storage:     NewStorage(storagePath),
		loginState:  StateUninitialized,
		isConnected: nil,
	}
}

func (m *Manager) HasConnection() bool {
	// Only try to connect once
	if m.isConnected != nil {
		return m.isConnected()
	}

	logger.Debug("Checking connection to %s", config.SERVER_ENDPOINT)
	resp, err := client.Get(fmt.Sprintf("%s/health", config.SERVER_ENDPOINT))
	if err != nil {
		logger.Error("Failed to ping server: %v", err)
		m.isConnected = func() bool { return false }
		return false
	}
	defer resp.Body.Close()

	ok := resp.StatusCode == http.StatusOK
	logger.Debug("Connection status: %v", ok)
	m.isConnected = func() bool { return ok }
	return ok
}

// Initialize loads saved state and attempts auto-login
func (m *Manager) Initialize() error {
	// Create default user
	m.currentUser = &User{
		Settings: GetDefaultSettings(),
	}

	// Load saved settings
	if settings, err := m.storage.LoadSettings(); err == nil {
		m.currentUser.Settings = settings
	}

	// Try auto-login if credentials exist
	if err := m.TryAutoLogin(); err == nil {
		return nil
	}

	m.loginState = StateOffline
	return nil
}

// GetUser returns the current user state
func (m *Manager) GetUser() *User {
	return m.currentUser
}

// GetSettings returns the current settings
func (m *Manager) GetSettings() *Settings {
	if m.currentUser == nil {
		return GetDefaultSettings()
	}
	return m.currentUser.Settings
}

// GetLoginState returns the current login state
func (m *Manager) GetLoginState() LoginState {
	return m.loginState
}

// GetSession returns the current session if logged in
func (m *Manager) GetSession() *Session {
	return m.session
}

func (m *Manager) Register(username, password string) error {
	if m.loginState == StateOnline {
		return fmt.Errorf("already logged in")
	}

	if err := m.storage.client.Register(username, password); err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	return nil
}

func (m *Manager) GetLeaderboard(song string, difficulty int) ([]Score, error) {
	if !m.HasConnection() {
		return []Score{}, nil
	}
	lb, err := m.storage.client.GetLeaderboard(song, fmt.Sprintf("%d", difficulty))
	if err != nil {
		return []Score{}, nil
	}
	return lb, nil
}

// Login attempts to log in with credentials
func (m *Manager) Login(username, password string, remember bool) error {
	m.loginState = StateLoggingIn
	m.rememberMe = remember

	tokens, err := m.storage.client.Login(username, password)
	if err != nil {
		m.loginState = StateOffline
		return fmt.Errorf("login failed: %w", err)
	}

	// Create new session
	m.session = &Session{
		Username:     username,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	// Get user
	user, err := m.storage.client.GetUser(tokens.AccessToken)
	if err != nil {
		m.loginState = StateOffline
		return fmt.Errorf("failed to get user: %w", err)
	}

	user.Settings = m.currentUser.Settings
	m.currentUser = user
	m.loginState = StateOnline

	// Save credentials if remember me enabled
	if remember {
		if err := m.storage.SaveCredentials(username, tokens.RefreshToken); err != nil {
			// Log but don't fail login
			logger.Error("Failed to save credentials: %v\n", err)
		}
	}

	// // Sync settings
	// if user.Settings != nil {
	// 	m.currentUser.Settings.MergeFrom(user.Settings)
	// 	if err := m.storage.SaveSettings(m.currentUser.Settings); err != nil {
	// 		fmt.Printf("Failed to save synced settings: %v\n", err)
	// 	}
	// }

	return nil
}

// TryAutoLogin attempts to login with stored credentials
func (m *Manager) TryAutoLogin() error {
	creds, err := m.storage.LoadCredentials()
	if err != nil || creds == nil {
		return fmt.Errorf("no stored credentials")
	}

	tokens, err := m.storage.client.Refresh(creds.RefreshToken)
	if err != nil {
		m.storage.ClearCredentials()
		return fmt.Errorf("stored credentials expired")
	}

	m.session = &Session{
		Username:     creds.Username,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	user, err := m.storage.client.GetUser(tokens.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	user.Settings = m.currentUser.Settings
	m.currentUser = user
	m.loginState = StateOnline
	if err := m.storage.SaveCredentials(user.Username, tokens.RefreshToken); err != nil {
		// Log but don't fail login
		logger.Error("Failed to save credentials: %v\n", err)
	}
	// Sync settings
	// if user.Settings != nil {
	// 	m.currentUser.Settings.MergeFrom(user.Settings)
	// 	if err := m.storage.SaveSettings(m.currentUser.Settings); err != nil {
	// 		fmt.Printf("Failed to save synced settings: %v\n", err)
	// 	}
	// }
	logger.Debug("Auto logged in as %s", m.currentUser.Username)

	return nil
}

// Logout ends the current session
func (m *Manager) Logout() {
	m.storage.ClearCredentials()

	if m.session != nil {
		m.session = nil
	}

	m.loginState = StateOffline
}

// SaveSettings persists current settings
func (m *Manager) SaveSettings() error {
	if m.currentUser == nil || m.currentUser.Settings == nil {
		return fmt.Errorf("no settings to save")
	}

	// Update modification time
	m.currentUser.Settings.LastModified = time.Now()

	// Save locally
	if err := m.storage.SaveSettings(m.currentUser.Settings); err != nil {
		return fmt.Errorf("failed to save settings: %w", err)
	}

	// // Sync to server if logged in
	// if m.loginState == StateOnline {
	// 	if err := m.storage.client.UpdateUser(m.session.AccessToken, m.currentUser.Settings); err != nil {
	// 		// Log but don't fail local save
	// 		fmt.Printf("Failed to sync settings to server: %v\n", err)
	// 	} else {
	// 		m.currentUser.Settings.LastSync = time.Now()
	// 	}
	// }

	return nil
}

func (m *Manager) AddScore(s *Score) error {
	if m.loginState != StateOnline {
		return nil
	}

	if err := m.storage.client.AddScore(m.session.AccessToken, s); err != nil {
		return fmt.Errorf("failed to add score: %w", err)
	}

	// Refetch user rank after
	user, err := m.storage.client.GetUser(m.session.AccessToken)
	if err != nil {
		logger.Debug("Failed to get user: %v", err)
		return fmt.Errorf("failed to get user: %w", err)
	}

	m.currentUser.Rank = user.Rank
	logger.Debug("Updated user rank from %d to %d", m.currentUser.Rank, user.Rank)
	return nil
}

func (m *Manager) GetScore(songHash string, difficulty int) (*Score, error) {
	if m.loginState != StateOnline {
		return nil, nil
	}

	s, err := m.storage.client.GetUserScores(m.session.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get score: %w", err)
	}
	for _, score := range s {
		if score.SongHash == songHash && score.Difficulty == difficulty {
			return &score, nil
		}
	}
	return nil, nil
}

// GetDefaultSettings returns default settings values
func GetDefaultSettings() *Settings {
	return &Settings{
		Version:            1,
		LastModified:       time.Now(),
		Locale:             "en-us",
		Fullscreen:         false,
		ScreenWidth:        1280,
		ScreenHeight:       720,
		RenderWidth:        1280,
		RenderHeight:       720,
		FixedRenderScale:   false,
		BGMVolume:          0.5,
		SFXVolume:          0.5,
		SongVolume:         0.7,
		LaneSpeed:          1.0,
		AudioOffset:        0,
		InputOffset:        0,
		NoteColorTheme:     "note.color.default",
		CenterNoteColor:    "#e6e600ff",
		CornerNoteColor:    "#e68200ff",
		DisableHoldNotes:   false,
		DisableHitEffects:  false,
		DisableLaneEffects: false,
		IsNewUser:          true,
		NoteWidth:          1.0,
		KeyConfig:          0,
	}
}

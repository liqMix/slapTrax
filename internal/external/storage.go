package external

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	settingsFilename = "settings.json"
	authFilename     = "auth.json"
)

// Storage handles persistent storage and server communication
type Storage struct {
	basePath string
	client   *APIClient
}

// StoredCredentials represents saved login information
type StoredCredentials struct {
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	SavedAt      time.Time `json:"saved_at"`
}

// NewStorage creates a new storage instance
func NewStorage(basePath string) *Storage {
	return &Storage{
		basePath: basePath,
		client:   NewAPIClient(),
	}
}

// SaveSettings persists settings to disk
func (s *Storage) SaveSettings(settings *Settings) error {
	path := filepath.Join(s.basePath, settingsFilename)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create settings directory: %w", err)
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}

	return nil
}

// LoadSettings reads settings from disk
func (s *Storage) LoadSettings() (*Settings, error) {
	path := filepath.Join(s.basePath, settingsFilename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read settings: %w", err)
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
	}

	return &settings, nil
}

// SaveCredentials stores login credentials
func (s *Storage) SaveCredentials(username, refreshToken string) error {
	path := filepath.Join(s.basePath, authFilename)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create auth directory: %w", err)
	}

	creds := StoredCredentials{
		Username:     username,
		RefreshToken: refreshToken,
		SavedAt:      time.Now(),
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials: %w", err)
	}

	return nil
}

// LoadCredentials reads stored credentials
func (s *Storage) LoadCredentials() (*StoredCredentials, error) {
	path := filepath.Join(s.basePath, authFilename)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials: %w", err)
	}

	var creds StoredCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
	}

	return &creds, nil
}

// ClearCredentials removes stored credentials
func (s *Storage) ClearCredentials() error {
	path := filepath.Join(s.basePath, authFilename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(path)
}

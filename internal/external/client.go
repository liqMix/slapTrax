package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/liqmix/ebiten-holiday-2024/internal/config"
)

// APIClient handles all server communication
type APIClient struct {
	baseURL string
	client  *http.Client
}

// NewAPIClient creates a new API client
func NewAPIClient() *APIClient {
	return &APIClient{
		baseURL: config.SERVER_ENDPOINT,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// Register creates a new user account
func (c *APIClient) Register(username, password string) error {
	body := map[string]interface{}{
		"username": username,
		"password": password,
	}

	resp, err := c.post("/register", body)
	if err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

// Login authenticates with username/password
func (c *APIClient) Login(username, password string) (*TokenPair, error) {
	body := map[string]string{
		"username": username,
		"password": password,
	}

	resp, err := c.post("/login", body)
	if err != nil {
		return nil, fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	var tokens TokenPair
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("failed to decode tokens: %w", err)
	}

	return &tokens, nil
}

// Refresh attempts to get new tokens using a refresh token
func (c *APIClient) Refresh(refreshToken string) (*TokenPair, error) {
	body := map[string]string{
		"refresh_token": refreshToken,
	}

	resp, err := c.post("/refresh", body)
	if err != nil {
		return nil, fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	var tokens TokenPair
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("failed to decode tokens: %w", err)
	}

	return &tokens, nil
}

// GetUser retrieves the user
func (c *APIClient) GetUser(accessToken string) (*User, error) {
	resp, err := c.authGet("/user", accessToken)
	if err != nil {
		return nil, fmt.Errorf("profile request failed: %w", err)
	}
	defer resp.Body.Close()

	var profile User
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("failed to decode profile: %w", err)
	}

	return &profile, nil
}

// UpdateUser updates the user
// func (c *APIClient) UpdateUser(accessToken string, settings *Settings) error {
// 	v, err := settings.Value()
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal settings: %w", err)
// 	}
// 	resp, err := c.authPost("/users", accessToken, v)
// 	if err != nil {
// 		return fmt.Errorf("profile update failed: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	return nil
// }

// Get user scores retrieves the scores for a user
func (c *APIClient) GetUserScores(accessToken string) ([]Score, error) {
	resp, err := c.authGet("/user/scores", accessToken)
	if err != nil {
		return nil, fmt.Errorf("scores request failed: %w", err)
	}
	defer resp.Body.Close()

	var scores []Score
	if err := json.NewDecoder(resp.Body).Decode(&scores); err != nil {
		return nil, fmt.Errorf("failed to decode scores: %w", err)
	}

	return scores, nil
}

// AddScore submits a new score
func (c *APIClient) AddScore(accessToken string, score *Score) error {
	resp, err := c.authPost("/scores", accessToken, score)
	if err != nil {
		return fmt.Errorf("score submission failed: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

// GetLeaderboard retrieves the leaderboard for a song
func (c *APIClient) GetLeaderboard(song, difficulty string) ([]Score, error) {
	resp, err := c.get("/scores/leaderboard?song=" + song + "&difficulty=" + difficulty)
	if err != nil {
		return nil, fmt.Errorf("leaderboard request failed: %w", err)
	}
	defer resp.Body.Close()

	var leaderboard []Score
	if err := json.NewDecoder(resp.Body).Decode(&leaderboard); err != nil {
		return nil, fmt.Errorf("failed to decode leaderboard: %w", err)
	}

	return leaderboard, nil
}

// get performs a GET request
func (c *APIClient) get(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("server returned %d: %s", resp.StatusCode, body)
	}

	return resp, nil
}

// post performs a POST request with JSON body
func (c *APIClient) post(path string, body interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("server returned %d: %s", resp.StatusCode, body)
	}

	return resp, nil
}

// authGet performs a GET request with auth header
func (c *APIClient) authGet(path, token string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("server returned %d: %s", resp.StatusCode, body)
	}

	return resp, nil
}

// authPost performs a POST request with auth header and JSON body
func (c *APIClient) authPost(path, token string, body interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("server returned %d: %s", resp.StatusCode, body)
	}

	return resp, nil
}

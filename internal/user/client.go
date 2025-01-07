package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/liqmix/ebiten-holiday-2024/internal/config"
)

var (
	apiURL = config.SERVER_ENDPOINT
	client = &http.Client{Timeout: 10 * time.Second}
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AutoLoginResponse struct {
	Found    bool   `json:"found"`
	Username string `json:"username,omitempty"`
	LastSeen string `json:"last_seen,omitempty"`
}

func CheckAutoLogin() (*AutoLoginResponse, error, int) {
	resp, err := client.Get(fmt.Sprintf("%s/auth/check-auto-login", apiURL))
	if err != nil {
		return nil, fmt.Errorf("failed to check auto-login: %w", err), 500
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode), resp.StatusCode
	}

	var response AutoLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err), 500
	}

	return &response, nil, http.StatusOK
}

func AttemptAutoLogin() error {
	resp, err := client.Post(fmt.Sprintf("%s/auth/auto-login", apiURL), "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to auto-login: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	var tokens TokenPair
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return fmt.Errorf("failed to decode tokens: %w", err)
	}

	currentUser.tokens = &tokens
	return nil
}

func Login(username, password string) error {
	credentials := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(credentials)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	resp, err := client.Post(fmt.Sprintf("%s/auth/login", apiURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	var tokens TokenPair
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return fmt.Errorf("failed to decode tokens: %w", err)
	}

	currentUser.tokens = &tokens
	return nil
}

func Refresh() error {
	if currentUser.tokens == nil {
		return fmt.Errorf("no refresh token available")
	}

	jsonData, err := json.Marshal(map[string]string{
		"refresh_token": currentUser.tokens.RefreshToken,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal refresh token: %w", err)
	}

	resp, err := client.Post(fmt.Sprintf("%s/auth/refresh", apiURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	var tokens TokenPair
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return fmt.Errorf("failed to decode tokens: %w", err)
	}

	currentUser.tokens = &tokens
	return nil
}

func makeAuthRequest(method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", apiURL, path), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if currentUser.tokens != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", currentUser.tokens.AccessToken))
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		// Try to refresh the token
		if err := Refresh(); err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}
		// Retry the request with new token
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", currentUser.tokens.AccessToken))
		return client.Do(req)
	}

	return resp, nil
}

func GetUserProfile() (*User, error) {
	if currentUser.tokens == nil {
		return nil, fmt.Errorf("not logged in")
	}

	resp, err := makeAuthRequest("GET", "/auth/profile", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	var profile User
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("failed to decode profile: %w", err)
	}

	return &profile, nil
}

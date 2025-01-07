package user

import "github.com/liqmix/ebiten-holiday-2024/internal/logger"

var currentUser *User
var S *Settings

type User struct {
	ID            uint      `json:"id"`
	Username      string    `json:"username"`
	settings      *Settings `json:"settings"`
	tokens        *TokenPair
	IsNewUser     bool
	IsGuest       bool
	HasServerUser bool
	BypassLogin   bool
}

func Init() error {
	currentUser = newUser()

	var err error
	currentUser.settings, err = LoadSettings()
	if err != nil {
		logger.Error("failed to load settings: %v", err)
		currentUser.settings = NewSettings()
		currentUser.IsNewUser = true
	}

	S = currentUser.settings

	// Check for auto-login possibility
	autoLoginResp, err, statusCode := CheckAutoLogin()
	if err != nil {
		if statusCode == 404 {
			currentUser.HasServerUser = false
		}

		return err
	}

	if autoLoginResp.Found {
		// Attempt auto-login
		if err := AttemptAutoLogin(); err != nil {
			return err
		}

		// Get user profile after successful auto-login
		profile, err := GetUserProfile()
		if err != nil {
			return err
		}

		currentUser.ID = profile.ID
		currentUser.Username = profile.Username
		currentUser.IsGuest = false
		currentUser.IsNewUser = false
		// Merge any remote settings with local settings if needed
		if profile.settings != nil {
			currentUser.settings.MergeFrom(profile.settings)
		}
	}

	return nil
}

func ShowLogin() bool {
	return false
	// return !currentUser.BypassLogin && currentUser.IsNewUser || !currentUser.HasServerUser
}

func newUser() *User {
	return &User{}
}

func Save() {
	if currentUser.IsGuest || currentUser == nil || currentUser.settings == nil {
		return
	}

	currentUser.settings.Save()
}

func IsLoggedIn() bool {
	return currentUser != nil && currentUser.tokens != nil
}

func GetCurrentUser() *User {
	return currentUser
}

func GetUsername() string {
	if currentUser == nil {
		return ""
	}
	return currentUser.Username
}

func Logout() {
	if currentUser != nil {
		currentUser.tokens = nil
		currentUser.ID = 0
		currentUser.Username = ""
	}
}

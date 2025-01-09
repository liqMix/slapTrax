package external

import (
	"net/http"
	"net/url"

	"time"

	"github.com/pkg/browser"
)

var (
	client         = &http.Client{Timeout: 10 * time.Second}
	M              = NewManager("storage")
	HasConnection  = M.HasConnection
	GetLoginState  = M.GetLoginState
	Logout         = M.Logout
	Login          = M.Login
	Register       = M.Register
	GetLeaderboard = M.GetLeaderboard
	AddScore       = M.AddScore
	PlayOffline    = M.PlayOffline
)

// Opens browser to URL
func OpenURL(path string) error {
	_, err := url.Parse(path)
	if err != nil {
		return err
	}

	return browser.OpenURL(path)
}

package external

import (
	"fmt"
	"net/http"
	"net/url"

	"time"

	"github.com/pkg/browser"

	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
)

var client = &http.Client{Timeout: 10 * time.Second}

// Opens browser to URL
func OpenURL(path string) error {
	_, err := url.Parse(path)
	if err != nil {
		return err
	}

	return browser.OpenURL(path)
}

func PingServer() bool {
	resp, err := client.Get(fmt.Sprintf("%s/health", config.SERVER_ENDPOINT))
	if err != nil {
		logger.Error("Failed to ping server: %v", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode != http.StatusOK
}

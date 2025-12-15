package browser

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/go-rod/rod/lib/launcher"
)

// Chrome represents a Chrome browser instance
type Chrome struct {
	launcher *launcher.Launcher
	url      string
	cancel   context.CancelFunc
}

// ChromeConfig holds Chrome configuration
type ChromeConfig struct {
	Headless       bool
	Width          int
	Height         int
	UserDataDir    string
	DisableWebSec  bool
}

// LaunchChrome launches a Chrome browser with remote debugging enabled
func LaunchChrome(ctx context.Context, cfg ChromeConfig) (*Chrome, error) {
	// Create user data directory if it doesn't exist
	if cfg.UserDataDir != "" {
		if err := os.MkdirAll(cfg.UserDataDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create user data dir: %w", err)
		}
	}

	// Create launcher
	l := launcher.New().
		Headless(cfg.Headless).
		Set("disable-blink-features", "AutomationControlled").
		Set("disable-features", "IsolateOrigins,site-per-process").
		Set("disable-infobars").
		Set("exclude-switches", "enable-automation").
		Set("disable-dev-shm-usage").
		Set("no-first-run").
		Set("no-default-browser-check").
		Set("disable-background-networking").
		Set("disable-background-timer-throttling").
		Set("disable-backgrounding-occluded-windows").
		Set("disable-breakpad").
		Set("disable-client-side-phishing-detection").
		Set("disable-default-apps").
		Set("disable-hang-monitor").
		Set("disable-popup-blocking").
		Set("disable-prompt-on-repost").
		Set("disable-sync").
		Set("metrics-recording-only").
		Set("no-service-autorun").
		Set("password-store", "basic").
		Set("use-mock-keychain")

	// Set window size
	if cfg.Width > 0 && cfg.Height > 0 {
		l = l.Set("window-size", fmt.Sprintf("%d,%d", cfg.Width, cfg.Height))
	}

	// Set user data directory for cookie persistence
	if cfg.UserDataDir != "" {
		l = l.UserDataDir(cfg.UserDataDir)
	}

	// Disable web security if requested (for testing only)
	if cfg.DisableWebSec {
		l = l.Set("disable-web-security")
	}

	// Create context with cancel
	launchCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Launch browser
	url, err := l.Context(launchCtx).Launch()
	if err != nil {
		return nil, fmt.Errorf("failed to launch Chrome: %w", err)
	}

	chrome := &Chrome{
		launcher: l,
		url:      url,
	}

	return chrome, nil
}

// URL returns the WebSocket URL for the browser
func (c *Chrome) URL() string {
	return c.url
}

// Close gracefully closes the Chrome browser
func (c *Chrome) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	
	if c.launcher != nil {
		c.launcher.Cleanup()
	}

	return nil
}

// KillChrome forcefully kills all Chrome processes (emergency cleanup)
func KillChrome() error {
	// This is a fallback for cleanup - not ideal but necessary for stuck instances
	cmd := exec.Command("pkill", "-f", "chrome")
	_ = cmd.Run() // Ignore errors as processes may not exist
	return nil
}

// GetUserDataDir returns the default user data directory path
func GetUserDataDir(dataDir string) string {
	return filepath.Join(dataDir, "chrome-user-data")
}

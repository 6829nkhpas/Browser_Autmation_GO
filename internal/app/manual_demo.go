package app

import (
	"fmt"
	"time"

	"github.com/nkh/linkedin-automation/internal/auth"
)

// runManualLoginDemo opens LinkedIn and waits for manual login
func (a *App) runManualLoginDemo() error {
	// Navigate to LinkedIn feed
	page := a.behavior.Page()
	err := a.behavior.Navigate("https://www.linkedin.com/feed/")
	if err != nil {
		return fmt.Errorf("failed to open LinkedIn: %w", err)
	}

	a.logger.Info("LinkedIn opened. Please login if not already logged in...")
	a.logger.Info("Waiting 60 seconds for you to complete login...")

	// Wait for manual login
	time.Sleep(60 * time.Second)

	// Check if logged in
	detector := auth.NewDetector(page)
	if !detector.IsLoggedIn() {
		a.logger.Info("Login not detected. Waiting another 30 seconds...")
		time.Sleep(30 * time.Second)

		if !detector.IsLoggedIn() {
			return fmt.Errorf("still not logged in after waiting")
		}
	}

	a.logger.Info("✅ Login detected!")
	a.logger.Info("Saving session cookies...")

	// Save cookies
	cookieStore := auth.NewCookieStore(a.config.Paths.CookieFile)
	if err := cookieStore.Save(page); err != nil {
		a.logger.Error(fmt.Sprintf("Failed to save cookies: %v", err))
	} else {
		a.logger.Info("✅ Cookies saved!")
	}

	a.logger.Info("")
	a.logger.Info("=== STARTING AUTOMATION DEMONSTRATION ===")
	time.Sleep(3 * time.Second)

	// Run the example automation flow
	return a.runExampleAutomation()
}

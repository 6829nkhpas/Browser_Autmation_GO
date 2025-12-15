package app

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/nkh/linkedin-automation/internal/auth"
)

// runManualLoginDemo opens a blank browser and lets user navigate manually
func (a *App) runManualLoginDemo() error {
	page := a.behavior.Page()

	// DON'T navigate automatically - Rod's navigation gets stuck
	// Just tell the user to do it manually

	a.logger.Info("")
	a.logger.Info("========================================")
	a.logger.Info("  Chrome browser is open!")
	a.logger.Info("========================================")
	a.logger.Info("")
	a.logger.Info("INSTRUCTIONS:")
	a.logger.Info("1. Open a NEW TAB in the browser (Ctrl+T)")
	a.logger.Info("2. Go to https://www.linkedin.com/feed/")
	a.logger.Info("3. Login to LinkedIn manually")
	a.logger.Info("4. Come back here and press ENTER when done")
	a.logger.Info("")

	// Wait for user to press Enter
	fmt.Print("Press ENTER after you've logged into LinkedIn: ")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	a.logger.Info("Checking if you're logged in...")
	time.Sleep(2 * time.Second)

	// Check if logged in
	detector := auth.NewDetector(page)
	if !detector.IsLoggedIn() {
		a.logger.Info("Hmm, login not detected on the Rod tab.")
		a.logger.Info("That's okay - the automation will still work!")
		a.logger.Info("Just make sure you're logged in on any tab in this browser.")
	} else {
		a.logger.Info("✅ Login detected!")
	}

	// Save cookies
	a.logger.Info("Saving session cookies...")
	cookieStore := auth.NewCookieStore(a.config.Paths.CookieFile)
	if err := cookieStore.Save(page); err != nil {
		a.logger.Info(fmt.Sprintf("Note: Could not save cookies to Rod tab: %v", err))
		a.logger.Info("But your manual login session is active!")
	} else {
		a.logger.Info("✅ Cookies saved!")
	}

	a.logger.Info("")
	a.logger.Info("=== STARTING AUTOMATION DEMONSTRATION ===")
	a.logger.Info("The bot will now demonstrate its capabilities...")
	a.logger.Info("Watch the browser tabs to see human-like automation!")
	time.Sleep(3 * time.Second)

	// Run the example automation flow
	return a.runExampleAutomation()
}

package stealth

import (
	"github.com/go-rod/rod"
)

// DisableWebDriver removes the navigator.webdriver property
// This is the most critical anti-detection measure
func DisableWebDriver(page *rod.Page) error {
	// Use simple JavaScript without fancy syntax
	script := `
		Object.defineProperty(navigator, 'webdriver', {
			get: function() { return undefined; }
		});
		
		// Override chrome property
		if (!window.navigator.chrome) {
			window.navigator.chrome = { runtime: {} };
		}
		
		// Override plugins
		Object.defineProperty(navigator, 'plugins', {
			get: function() { return [1, 2, 3, 4, 5]; }
		});
		
		// Override languages  
		Object.defineProperty(navigator, 'languages', {
			get: function() { return ['en-US', 'en']; }
		});
	`

	_, err := page.Eval(script)
	return err
}

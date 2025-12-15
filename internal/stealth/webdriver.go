package stealth

import (
	"github.com/go-rod/rod"
)

// DisableWebDriver removes the navigator.webdriver property
// This is the most critical anti-detection measure
func DisableWebDriver(page *rod.Page) error {
	script := `
		Object.defineProperty(navigator, 'webdriver', {
			get: function() { return undefined; }
		});
	`

	_, err := page.Eval(script)
	return err
}

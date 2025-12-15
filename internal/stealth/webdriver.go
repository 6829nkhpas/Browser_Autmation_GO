package stealth

import (
	"github.com/go-rod/rod"
)

// DisableWebDriver removes the navigator.webdriver property
// This is the most critical anti-detection measure
func DisableWebDriver(page *rod.Page) error {
	script := `
		// Remove webdriver property
		Object.defineProperty(navigator, 'webdriver', {
			get: () => undefined
		});

		// Override automation-related properties
		window.navigator.chrome = {
			runtime: {},
		};

		// Override permissions
		const originalQuery = window.navigator.permissions.query;
		window.navigator.permissions.query = (parameters) => (
			parameters.name === 'notifications' ?
				Promise.resolve({ state: Notification.permission }) :
				originalQuery(parameters)
		);

		// Override plugins to appear more realistic
		Object.defineProperty(navigator, 'plugins', {
			get: () => [1, 2, 3, 4, 5]
		});

		// Override languages
		Object.defineProperty(navigator, 'languages', {
			get: () => ['en-US', 'en']
		});
	`

	_, err := page.Eval(script)
	return err
}

package stealth

import (
	"github.com/go-rod/rod"
)

// SpoofLocale sets the timezone and language for the browser
func SpoofLocale(page *rod.Page, timezone, language string) error {
	if timezone == "" {
		timezone = "America/New_York"
	}
	if language == "" {
		language = "en-US"
	}

	script := `
		// Override timezone
		Object.defineProperty(Intl.DateTimeFormat.prototype, 'resolvedOptions', {
			value: function() {
				return {
					locale: '` + language + `',
					calendar: 'gregory',
					numberingSystem: 'latn',
					timeZone: '` + timezone + `',
					year: 'numeric',
					month: '2-digit',
					day: '2-digit'
				};
			}
		});

		// Override language
		Object.defineProperty(navigator, 'language', {
			get: () => '` + language + `'
		});

		Object.defineProperty(navigator, 'languages', {
			get: () => ['` + language + `', '` + language[:2] + `']
		});
	`

	_, err := page.Eval(script)
	return err
}

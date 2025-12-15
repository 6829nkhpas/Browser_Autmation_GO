package stealth

import (
	"github.com/go-rod/rod"
)

// OverridePermissions overrides the Permissions API to avoid detection
func OverridePermissions(page *rod.Page) error {
	script := `
		var originalQuery = window.navigator.permissions.query;
		window.navigator.permissions.query = function(params) {
			return Promise.resolve({ state: 'prompt', onchange: null });
		};
	`

	_, err := page.Eval(script)
	return err
}

package stealth

import (
	"github.com/go-rod/rod"
)

// OverridePermissions overrides the Permissions API to avoid detection
func OverridePermissions(page *rod.Page) error {
	script := `
		// Override permissions query
		const originalQuery = window.navigator.permissions.query;
		
		window.navigator.permissions.query = function(parameters) {
			// Return sensible defaults for common permissions
			if (parameters.name === 'notifications') {
				return Promise.resolve({
					state: 'prompt',
					onchange: null
				});
			}
			
			if (parameters.name === 'geolocation') {
				return Promise.resolve({
					state: 'prompt',
					onchange: null
				});
			}
			
			if (parameters.name === 'camera' || parameters.name === 'microphone') {
				return Promise.resolve({
					state: 'prompt',
					onchange: null
				});
			}
			
			// For other permissions, use original query
			return originalQuery(parameters);
		};
	`

	return page.Eval(script).Error
}

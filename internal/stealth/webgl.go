package stealth

import (
	"crypto/rand"
	"math/big"

	"github.com/go-rod/rod"
)

var vendors = []string{
	"Intel Inc.",
	"NVIDIA Corporation",
	"AMD",
	"Intel",
}

var renderers = []string{
	"Intel Iris OpenGL Engine",
	"NVIDIA GeForce GTX 1060",
	"AMD Radeon RX 580",
	"Intel(R) UHD Graphics 630",
	"ANGLE (Intel, Intel(R) UHD Graphics 630 Direct3D11 vs_5_0 ps_5_0)",
}

// RandomizeWebGL randomizes WebGL fingerprint
func RandomizeWebGL(page *rod.Page) error {
	vendor := getRandomItem(vendors)
	renderer := getRandomItem(renderers)

	script := `
		(function() {
			const getParameter = WebGLRenderingContext.prototype.getParameter;
			
			WebGLRenderingContext.prototype.getParameter = function(parameter) {
				// UNMASKED_VENDOR_WEBGL
				if (parameter === 37445) {
					return '` + vendor + `';
				}
				// UNMASKED_RENDERER_WEBGL
				if (parameter === 37446) {
					return '` + renderer + `';
				}
				return getParameter.apply(this, arguments);
			};
		})();
	`

	_, err := page.Eval(script)
	return err
}

func getRandomItem(items []string) string {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(items))))
	if err != nil {
		return items[0]
	}
	return items[n.Int64()]
}

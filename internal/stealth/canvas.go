package stealth

import (
	"crypto/rand"
	"math/big"

	"github.com/go-rod/rod"
)

// RandomizeCanvas adds subtle noise to canvas fingerprinting
// This makes each session unique while staying consistent within the session
func RandomizeCanvas(page *rod.Page) error {
	// Generate a small random noise value for this session
	noise := getCanvasNoise()

	script := `(function() {
		const originalToDataURL = HTMLCanvasElement.prototype.toDataURL;
		const originalGetImageData = CanvasRenderingContext2D.prototype.getImageData;
		
		const noise = ` + noise + `;
		
		// Override toDataURL
		HTMLCanvasElement.prototype.toDataURL = function() {
			const context = this.getContext('2d');
			const imageData = context.getImageData(0, 0, this.width, this.height);
			
			// Add subtle noise to image data
			for (let i = 0; i < imageData.data.length; i += 4) {
				imageData.data[i] = imageData.data[i] + noise;
				imageData.data[i + 1] = imageData.data[i + 1] + noise;
				imageData.data[i + 2] = imageData.data[i + 2] + noise;
			}
			
			context.putImageData(imageData, 0, 0);
			return originalToDataURL.apply(this, arguments);
		};
		
		// Override getImageData to add consistent noise
		CanvasRenderingContext2D.prototype.getImageData = function() {
			const imageData = originalGetImageData.apply(this, arguments);
			for (let i = 0; i < imageData.data.length; i += 4) {
				imageData.data[i] = imageData.data[i] + noise;
				imageData.data[i + 1] = imageData.data[i + 1] + noise;
				imageData.data[i + 2] = imageData.data[i + 2] + noise;
			}
			return imageData;
		};
	})()`

	_, err := page.Eval(script)
	return err
}

func getCanvasNoise() string {
	// Generate random noise between -2 and 2
	n, err := rand.Int(rand.Reader, big.NewInt(5))
	if err != nil {
		return "1"
	}
	return string(rune('0' + int(n.Int64()) - 2))
}

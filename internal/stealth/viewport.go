package stealth

import (
	"crypto/rand"
	"math/big"
)

// ViewportSize represents a viewport dimension
type ViewportSize struct {
	Width  int
	Height int
}

var commonViewports = []ViewportSize{
	{Width: 1920, Height: 1080},
	{Width: 1366, Height: 768},
	{Width: 1536, Height: 864},
	{Width: 1440, Height: 900},
	{Width: 1280, Height: 720},
	{Width: 1600, Height: 900},
	{Width: 2560, Height: 1440},
}

// GetRandomViewport returns a random common viewport size
func GetRandomViewport() ViewportSize {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(commonViewports))))
	if err != nil {
		// Fallback to most common viewport
		return commonViewports[0]
	}
	return commonViewports[n.Int64()]
}

// GetViewport returns a specific viewport or a random one
func GetViewport(width, height int, randomize bool) ViewportSize {
	if width > 0 && height > 0 {
		return ViewportSize{Width: width, Height: height}
	}
	if randomize {
		return GetRandomViewport()
	}
	return commonViewports[0]
}

// AddVariation adds small random variation to viewport size
// This makes each session unique while staying within realistic bounds
func (v ViewportSize) AddVariation() ViewportSize {
	// Add variation of +/- 50 pixels
	widthVar := getRandomVariation(50)
	heightVar := getRandomVariation(50)

	return ViewportSize{
		Width:  v.Width + widthVar,
		Height: v.Height + heightVar,
	}
}

func getRandomVariation(maxVariation int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(maxVariation*2+1)))
	if err != nil {
		return 0
	}
	return int(n.Int64()) - maxVariation
}

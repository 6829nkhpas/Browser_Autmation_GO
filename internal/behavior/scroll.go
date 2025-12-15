package behavior

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
)

// ScrollNatural scrolls the page naturally with acceleration and deceleration
func ScrollNatural(page *rod.Page, totalPixels int) error {
	if totalPixels == 0 {
		return nil
	}

	// Determine scroll direction
	direction := 1
	if totalPixels < 0 {
		direction = -1
		totalPixels = -totalPixels
	}

	// Number of scroll steps
	steps := 10 + GetRandomInRange(0, 5)
	scrolled := 0

	for i := 0; i < steps && scrolled < totalPixels; i++ {
		// Calculate scroll amount for this step using easing function
		t := float64(i) / float64(steps)
		easedProgress := easeInOutCubic(t)

		nextScrolled := int(easedProgress * float64(totalPixels))
		delta := (nextScrolled - scrolled) * direction

		// Perform scroll
		err := page.Mouse.Scroll(0, float64(delta), steps)
		if err != nil {
			return fmt.Errorf("scroll failed: %w", err)
		}

		scrolled = nextScrolled

		// Variable delay between scroll steps
		delay := 30 + GetRandomInRange(0, 40)
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	return nil
}

// ScrollToElement scrolls to make an element visible
func ScrollToElement(page *rod.Page, selector string) error {
	elem, err := page.Element(selector)
	if err != nil {
		return fmt.Errorf("element not found: %w", err)
	}

	if err := elem.ScrollIntoView(); err != nil {
		return fmt.Errorf("scroll into view failed: %w", err)
	}

	// Brief pause after scrolling
	WaitHuman(300, 700)

	return nil
}

// RandomScroll performs random scrolling behavior (exploring page)
func RandomScroll(page *rod.Page) error {
	// Random number of scroll actions
	scrolls := GetRandomInRange(2, 6)

	for i := 0; i < scrolls; i++ {
		// Random scroll amount
		pixels := GetRandomInRange(200, 600)

		// Sometimes scroll up instead of down
		if GetRandomInRange(0, 100) < 20 {
			pixels = -pixels
		}

		if err := ScrollNatural(page, pixels); err != nil {
			return err
		}

		// Pause as if reading
		WaitHuman(1000, 3000)
	}

	return nil
}

// easeInOutCubic provides smooth acceleration and deceleration
// This makes scrolling feel natural
func easeInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	p := 2*t - 2
	return 1 + 0.5*p*p*p
}

// easeOutQuad provides deceleration effect
func easeOutQuad(t float64) float64 {
	return t * (2 - t)
}

// ScrollPageSection scrolls through a section with reading pauses
func ScrollPageSection(page *rod.Page, sectionHeight int) error {
	scrolled := 0

	for scrolled < sectionHeight {
		// Scroll a viewport-sized chunk
		chunk := GetRandomInRange(300, 500)
		if scrolled+chunk > sectionHeight {
			chunk = sectionHeight - scrolled
		}

		if err := ScrollNatural(page, chunk); err != nil {
			return err
		}

		scrolled += chunk

		// Reading pause
		WaitHuman(800, 2000)

		// Occasionally scroll back up a bit (re-reading)
		if shouldScrollBack() {
			backScroll := -GetRandomInRange(50, 150)
			if err := ScrollNatural(page, backScroll); err != nil {
				return err
			}
			WaitHuman(500, 1000)
		}
	}

	return nil
}

func shouldScrollBack() bool {
	return GetRandomInRange(0, 100) < 15 // 15% chance
}

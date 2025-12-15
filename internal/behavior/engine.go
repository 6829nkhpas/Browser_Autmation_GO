package behavior

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// Engine implements human-like browser behavior
// ALL browser interactions must go through this engine
type Engine struct {
	page *rod.Page
	ctx  context.Context
}

// New creates a new behavior engine
func New(ctx context.Context, page *rod.Page) *Engine {
	return &Engine{
		page: page,
		ctx:  ctx,
	}
}

// Navigate navigates to a URL with realistic delays
func (e *Engine) Navigate(url string) error {
	// Before navigation, add thinking delay
	WaitHuman(500, 1500)

	err := e.page.Navigate(url)
	if err != nil {
		return fmt.Errorf("navigation failed: %w", err)
	}

	// Wait for page load
	if err := e.page.WaitLoad(); err != nil {
		return fmt.Errorf("page load wait failed: %w", err)
	}

	// Add post-navigation delay (page processing time)
	WaitHuman(1000, 2000)

	return nil
}

// Click performs a human-like click on an element
// This includes: hover → pause → click
func (e *Engine) Click(selector string) error {
	// Find the element
	elem, err := e.page.Element(selector)
	if err != nil {
		return fmt.Errorf("element not found %s: %w", selector, err)
	}

	// Scroll element into view if needed
	if err := elem.ScrollIntoView(); err != nil {
		return fmt.Errorf("scroll into view failed: %w", err)
	}

	// Small delay after scrolling
	WaitHuman(300, 700)

	// Hover before clicking (critical for human behavior)
	if err := e.Hover(selector); err != nil {
		return fmt.Errorf("hover failed: %w", err)
	}

	// Thinking pause before click
	WaitHuman(200, 500)

	// Perform click
	if err := elem.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("click failed: %w", err)
	}

	// Post-click delay
	WaitHuman(500, 1000)

	return nil
}

// Type types text into an input field with human-like behavior
func (e *Engine) Type(selector, text string) error {
	// Find the element
	elem, err := e.page.Element(selector)
	if err != nil {
		return fmt.Errorf("element not found %s: %w", selector, err)
	}

	// Click on the element first to focus
	if err := elem.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("click to focus failed: %w", err)
	}

	WaitHuman(200, 400)

	// Type with human-like behavior
	if err := TypeHumanLike(elem, text); err != nil {
		return fmt.Errorf("typing failed: %w", err)
	}

	return nil
}

// Scroll scrolls the page by a certain number of pixels with natural acceleration
func (e *Engine) Scroll(pixels int) error {
	return ScrollNatural(e.page, pixels)
}

// Hover moves the mouse to an element with Bézier curve movement
func (e *Engine) Hover(selector string) error {
	elem, err := e.page.Element(selector)
	if err != nil {
		return fmt.Errorf("element not found %s: %w", selector, err)
	}

	// Get element position
	shape, err := elem.Shape()
	if err != nil {
		return fmt.Errorf("failed to get element shape: %w", err)
	}

	// Get bounding box from shape
	if len(shape.Quads) == 0 {
		return fmt.Errorf("element has no quads")
	}

	// Calculate center from first quad
	quad := shape.Quads[0]
	targetX := (quad[0] + quad[2]) / 2
	targetY := (quad[1] + quad[5]) / 2

	// Move mouse with Bézier curve
	if err := MoveMouse(e.page, targetX, targetY); err != nil {
		return fmt.Errorf("mouse movement failed: %w", err)
	}

	// Small pause after hover
	WaitHuman(100, 300)

	return nil
}

// WaitHuman adds a random human-like delay
func (e *Engine) WaitHuman(minMs, maxMs int) {
	WaitHuman(minMs, maxMs)
}

// ScrollToBottom scrolls to the bottom of the page naturally
func (e *Engine) ScrollToBottom() error {
	// Get page height
	height, err := e.page.Eval(`() => document.body.scrollHeight`)
	if err != nil {
		return fmt.Errorf("failed to get page height: %w", err)
	}

	totalHeight := int(height.Value.Num())
	scrolled := 0

	// Scroll in chunks
	for scrolled < totalHeight {
		chunkSize := GetRandomInRange(300, 600)
		if scrolled+chunkSize > totalHeight {
			chunkSize = totalHeight - scrolled
		}

		if err := e.Scroll(chunkSize); err != nil {
			return err
		}

		scrolled += chunkSize

		// Random pause while scrolling
		WaitHuman(500, 1500)
	}

	return nil
}

// WaitForElement waits for an element to appear with timeout
func (e *Engine) WaitForElement(selector string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(e.ctx, timeout)
	defer cancel()

	_, err := e.page.Context(ctx).Element(selector)
	if err != nil {
		return fmt.Errorf("element %s did not appear within %v: %w", selector, timeout, err)
	}

	return nil
}

// HasElement checks if an element exists on the page
func (e *Engine) HasElement(selector string) bool {
	has, _, err := e.page.Has(selector)
	return err == nil && has
}

// Page returns the underlying rod.Page for direct access when needed
func (e *Engine) Page() *rod.Page {
	return e.page
}

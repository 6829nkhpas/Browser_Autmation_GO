package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// Rod represents a Rod browser client
type Rod struct {
	browser *rod.Browser
	page    *rod.Page
}

// RodConfig holds Rod configuration
type RodConfig struct {
	ChromeURL string
	Timeout   time.Duration
}

// NewRod creates a new Rod browser client connected to Chrome
func NewRod(ctx context.Context, cfg RodConfig) (*Rod, error) {
	if cfg.ChromeURL == "" {
		return nil, fmt.Errorf("ChromeURL is required")
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	// Create browser instance
	browser := rod.New().
		ControlURL(cfg.ChromeURL).
		Context(ctx)

	// Connect to browser with timeout
	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to browser: %w", err)
	}

	// Test connection by getting version
	_, err := browser.Version()
	if err != nil {
		return nil, fmt.Errorf("failed to get browser version: %w", err)
	}

	return &Rod{
		browser: browser,
	}, nil
}

// Browser returns the underlying Rod browser instance
func (r *Rod) Browser() *rod.Browser {
	return r.browser
}

// NewPage creates a new page in the browser
func (r *Rod) NewPage(ctx context.Context) (*rod.Page, error) {
	// Create a new page using Rod's default method
	// This ensures all browser features work normally
	page, err := r.browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	// Set a reasonable timeout
	page = page.Timeout(30 * time.Second)

	// Store reference to current page
	r.page = page

	return page, nil
}

// CurrentPage returns the current active page
func (r *Rod) CurrentPage() *rod.Page {
	return r.page
}

// Close closes the Rod browser connection
func (r *Rod) Close() error {
	if r.page != nil {
		if err := r.page.Close(); err != nil {
			// Log error but don't fail
			fmt.Printf("warning: failed to close page: %v\n", err)
		}
	}

	if r.browser != nil {
		if err := r.browser.Close(); err != nil {
			return fmt.Errorf("failed to close browser: %w", err)
		}
	}

	return nil
}

// WaitLoad waits for the page to be fully loaded
func (r *Rod) WaitLoad() error {
	if r.page == nil {
		return fmt.Errorf("no active page")
	}

	return r.page.WaitLoad()
}

// WaitIdle waits for the page to be idle (no network activity)
func (r *Rod) WaitIdle(timeout time.Duration) error {
	if r.page == nil {
		return fmt.Errorf("no active page")
	}

	return r.page.WaitIdle(timeout)
}

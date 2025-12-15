package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/nkh/linkedin-automation/internal/behavior"
	"github.com/nkh/linkedin-automation/internal/linkedin"
)

// Authenticator handles LinkedIn authentication
type Authenticator struct {
	email       string
	password    string
	cookieStore *CookieStore
	behaviorEng *behavior.Engine
	detector    *Detector
}

// Config holds authentication configuration
type Config struct {
	Email      string
	Password   string
	CookieFile string
}

// New creates a new authenticator
func New(cfg Config, page *rod.Page, ctx context.Context) *Authenticator {
	return &Authenticator{
		email:       cfg.Email,
		password:    cfg.Password,
		cookieStore: NewCookieStore(cfg.CookieFile),
		behaviorEng: behavior.New(ctx, page),
		detector:    NewDetector(page),
	}
}

// Login performs the login flow
func (a *Authenticator) Login(ctx context.Context) error {
	// Check if already logged in via cookies
	if a.cookieStore.Exists() {
		if err := a.cookieStore.Load(a.behaviorEng.Page()); err != nil {
			return fmt.Errorf("failed to load cookies: %w", err)
		}

		// Navigate to LinkedIn home to verify session
		if err := a.behaviorEng.Navigate("https://www.linkedin.com/feed/"); err != nil {
			return fmt.Errorf("failed to navigate to LinkedIn: %w", err)
		}

		// Check if we're logged in
		if a.detector.IsLoggedIn() {
			return nil // Successfully restored session
		}

		// Cookies expired or invalid, delete them
		_ = a.cookieStore.Delete()
	}

	// Perform fresh login
	return a.performLogin(ctx)
}

// performLogin performs the actual login process
func (a *Authenticator) performLogin(ctx context.Context) error {
	// Navigate to login page
	if err := a.behaviorEng.Navigate(linkedin.LoginURL); err != nil {
		return fmt.Errorf("failed to navigate to login page: %w", err)
	}

	// Wait for page to load - give it more time
	page := a.behaviorEng.Page()
	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("page load timeout: %w", err)
	}

	// Additional wait for JavaScript to execute
	time.Sleep(3 * time.Second)

	// Check if already logged in by looking for feed or profile
	detector := NewDetector(page) // Assuming NewDetector signature is (page *rod.Page)
	if detector.IsLoggedIn() {
		return nil
	}

	// Check current URL to see if we're actually on login page
	currentURL := page.MustInfo().URL
	fmt.Printf("Current URL after navigation: %s\n", currentURL)

	// Wait for email input field to appear - this confirms page loaded
	_, err := page.Timeout(10 * time.Second).Element("#username")
	if err != nil {
		// Page didn't load properly, try to get page content for debugging
		html, _ := page.HTML()
		htmlLen := len(html)
		previewLen := 500
		if htmlLen < previewLen {
			previewLen = htmlLen
		}
		fmt.Printf("Page HTML length: %d\n", htmlLen)
		fmt.Printf("First %d chars: %s\n", previewLen, html[:previewLen])
		return fmt.Errorf("login page did not load properly - email field not found: %w", err)
	}

	fmt.Println("Login page loaded successfully, email field found")

	// Check for security challenge before login
	challenge, challengeErr := a.detector.DetectChallenge()
	if challengeErr != nil {
		return fmt.Errorf("failed to detect challenge: %w", challengeErr)
	}

	if challenge != ChallengeNone {
		return fmt.Errorf("security challenge detected before login: %s", challenge)
	}

	// Fill email
	emailSelector := "input[name='session_key'], input#username"
	if !a.behaviorEng.HasElement(emailSelector) {
		return fmt.Errorf("email input not found")
	}

	if err := a.behaviorEng.Type(emailSelector, a.email); err != nil {
		return fmt.Errorf("failed to type email: %w", err)
	}

	a.behaviorEng.WaitHuman(500, 1000)

	// Fill password
	passwordSelector := "input[name='session_password'], input#password"
	if !a.behaviorEng.HasElement(passwordSelector) {
		return fmt.Errorf("password input not found")
	}

	// Type password (no typos for passwords)
	pageForPassword := a.behaviorEng.Page()
	elem, err := pageForPassword.Element(passwordSelector)
	if err != nil {
		return fmt.Errorf("failed to find password field: %w", err)
	}

	if err := behavior.TypePassword(elem, a.password); err != nil {
		return fmt.Errorf("failed to type password: %w", err)
	}

	a.behaviorEng.WaitHuman(800, 1500)

	// Click sign in button
	signInSelector := "button[type='submit'], button[data-litms-control-urn*='login']"
	if !a.behaviorEng.HasElement(signInSelector) {
		return fmt.Errorf("sign in button not found")
	}

	if err := a.behaviorEng.Click(signInSelector); err != nil {
		return fmt.Errorf("failed to click sign in: %w", err)
	}

	// Wait for navigation
	a.behaviorEng.WaitHuman(3000, 5000)

	// Check for security challenges
	challenge2, challengeErr2 := a.detector.DetectChallenge()
	if challengeErr2 != nil {
		return fmt.Errorf("failed to detect post-login challenge: %w", challengeErr2)
	}

	if challenge2 != ChallengeNone {
		return fmt.Errorf("security challenge detected: %s - manual intervention required", challenge2)
	}

	// Verify login success
	if !a.detector.IsLoggedIn() {
		// Check if still on login page (indicates failure)
		if a.detector.IsLoginPage() {
			return fmt.Errorf("login failed - still on login page (check credentials)")
		}
		return fmt.Errorf("login state unclear - not on login page but not detected as logged in")
	}

	// Save cookies for next time
	if err := a.cookieStore.Save(page); err != nil {
		// Log warning but don't fail - we're logged in
		fmt.Printf("warning: failed to save cookies: %v\n", err)
	}

	return nil
}

// Logout logs out and clears cookies
func (a *Authenticator) Logout() error {
	// Navigate to logout URL would go here
	// For now, just clear cookies
	return a.cookieStore.Delete()
}

// IsLoggedIn checks if currently logged in
func (a *Authenticator) IsLoggedIn() bool {
	return a.detector.IsLoggedIn()
}

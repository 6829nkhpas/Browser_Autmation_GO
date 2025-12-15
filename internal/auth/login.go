package auth

import (
	"context"
	"fmt"

	"github.com/go-rod/rod"
	"github.com/nkh/linkedin-automation/internal/behavior"
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
	if err := a.behaviorEng.Navigate("https://www.linkedin.com/login"); err != nil {
		return fmt.Errorf("failed to navigate to login page: %w", err)
	}

	// Wait for page to load
	a.behaviorEng.WaitHuman(1000, 2000)

	// Check for security challenge before login
	challenge, err := a.detector.DetectChallenge()
	if err != nil {
		return fmt.Errorf("failed to detect challenge: %w", err)
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
	page := a.behaviorEng.Page()
	elem, err := page.Element(passwordSelector)
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
	challenge, err = a.detector.DetectChallenge()
	if err != nil {
		return fmt.Errorf("failed to detect post-login challenge: %w", err)
	}

	if challenge != ChallengeNone {
		return fmt.Errorf("security challenge detected: %s - manual intervention required", challenge)
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

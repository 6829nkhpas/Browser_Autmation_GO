package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
)

// Detector detects security challenges and checkpoints
type Detector struct {
	page *rod.Page
}

// NewDetector creates a new security challenge detector
func NewDetector(page *rod.Page) *Detector {
	return &Detector{
		page: page,
	}
}

// ChallengeType represents the type of security challenge
type ChallengeType string

const (
	ChallengeNone       ChallengeType = "none"
	ChallengeCaptcha    ChallengeType = "captcha"
	ChallengeOTP        ChallengeType = "otp"
	ChallengeCheckpoint ChallengeType = "checkpoint"
	ChallengeUnknown    ChallengeType = "unknown"
)

// DetectChallenge detects if there's a security challenge on the current page
func (d *Detector) DetectChallenge() (ChallengeType, error) {
	// Check for CAPTCHA
	if has := d.hasCaptcha(); has {
		return ChallengeCaptcha, nil
	}

	// Check for OTP/2FA
	if has := d.hasOTP(); has {
		return ChallengeOTP, nil
	}

	// Check for security checkpoint
	if has := d.hasCheckpoint(); has {
		return ChallengeCheckpoint, nil
	}

	return ChallengeNone, nil
}

// hasCaptcha checks for CAPTCHA presence
func (d *Detector) hasCaptcha() bool {
	captchaSelectors := []string{
		"iframe[src*='recaptcha']",
		"#recaptcha",
		".g-recaptcha",
		"iframe[src*='captcha']",
		"[data-testid='captcha']",
		".captcha",
	}

	for _, selector := range captchaSelectors {
		has, _, _ := d.page.Has(selector)
		if has {
			return true
		}
	}

	return false
}

// hasOTP checks for OTP/2FA prompt
func (d *Detector) hasOTP() bool {
	otpSelectors := []string{
		"input[name='pin']",
		"input[type='tel'][maxlength='6']",
		"input[placeholder*='code']",
		"input[placeholder*='verification']",
		"[data-testid='otp']",
		"[data-testid='verification-code']",
	}

	for _, selector := range otpSelectors {
		has, _, _ := d.page.Has(selector)
		if has {
			return true
		}
	}

	// Check for text content indicating OTP
	text, _ := d.page.HTML()
	otpKeywords := []string{
		"enter the code",
		"verification code",
		"two-factor",
		"2-factor",
		"security code",
		"authenticate your account",
	}

	for _, keyword := range otpKeywords {
		if contains(text, keyword) {
			return true
		}
	}

	return false
}

// hasCheckpoint checks for LinkedIn security checkpoint
func (d *Detector) hasCheckpoint() bool {
	checkpointSelectors := []string{
		"[data-testid='checkpoint']",
		".checkpoint-challenge",
		"#checkpoint",
	}

	for _, selector := range checkpointSelectors {
		has, _, _ := d.page.Has(selector)
		if has {
			return true
		}
	}

	// Check URL for checkpoint
	info, err := d.page.Info()
	if err == nil && contains(info.URL, "checkpoint") {
		return true
	}

	// Check for text indicating checkpoint
	text, _ := d.page.HTML()
	checkpointKeywords := []string{
		"security challenge",
		"verify your identity",
		"unusual activity",
		"confirm it's you",
	}

	for _, keyword := range checkpointKeywords {
		if contains(text, keyword) {
			return true
		}
	}

	return false
}

// WaitForChallengeResolution waits for user to resolve challenge manually
// This should only be used in non-headless mode for manual intervention
func (d *Detector) WaitForChallengeResolution(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		challenge, err := d.DetectChallenge()
		if err != nil {
			return err
		}

		if challenge == ChallengeNone {
			return nil // Challenge resolved
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
			// Continue checking
		}
	}

	return fmt.Errorf("challenge not resolved within timeout")
}

// IsLoginPage checks if current page is the login page
func (d *Detector) IsLoginPage() bool {
	info, err := d.page.Info()
	if err != nil {
		return false
	}

	// Check URL
	if contains(info.URL, "/login") || contains(info.URL, "/uas/login") {
		return true
	}

	// Check for login form elements
	loginSelectors := []string{
		"input[name='session_key']",
		"input[type='email'][name*='username']",
		"input[type='password']",
		"#username",
		"#password",
	}

	for _, selector := range loginSelectors {
		has, _, _ := d.page.Has(selector)
		if has {
			return true
		}
	}

	return false
}

// IsLoggedIn checks if user is logged in
func (d *Detector) IsLoggedIn() bool {
	// If on login page, definitely not logged in
	if d.IsLoginPage() {
		return false
	}

	// Check for elements that only appear when logged in
	loggedInSelectors := []string{
		"[data-testid='global-nav']",
		".global-nav",
		"#global-nav",
		".feed-identity-module",
	}

	for _, selector := range loggedInSelectors {
		has, _, _ := d.page.Has(selector)
		if has {
			return true
		}
	}

	// Check URL - logged in users typically see feed or profile
	info, err := d.page.Info()
	if err == nil {
		if contains(info.URL, "/feed") || contains(info.URL, "/in/") {
			return true
		}
	}

	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// CookieStore manages cookie persistence
type CookieStore struct {
	filePath string
}

// Cookie represents a browser cookie for storage
type Cookie struct {
	Name     string    `json:"name"`
	Value    string    `json:"value"`
	Domain   string    `json:"domain"`
	Path     string    `json:"path"`
	Expires  time.Time `json:"expires"`
	HTTPOnly bool      `json:"httpOnly"`
	Secure   bool      `json:"secure"`
	SameSite string    `json:"sameSite"`
}

// NewCookieStore creates a new cookie store
func NewCookieStore(filePath string) *CookieStore {
	return &CookieStore{
		filePath: filePath,
	}
}

// Save saves cookies to file
func (cs *CookieStore) Save(page *rod.Page) error {
	// Get all cookies from the browser
	cookies, err := page.Cookies([]string{})
	if err != nil {
		return fmt.Errorf("failed to get cookies: %w", err)
	}

	// Convert to our Cookie struct
	storeCookies := make([]Cookie, 0, len(cookies))
	for _, c := range cookies {
		storeCookies = append(storeCookies, Cookie{
			Name:     c.Name,
			Value:    c.Value,
			Domain:   c.Domain,
			Path:     c.Path,
			Expires:  time.Unix(int64(c.Expires), 0),
			HTTPOnly: c.HTTPOnly,
			Secure:   c.Secure,
			SameSite: string(c.SameSite),
		})
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(cs.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(storeCookies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cookies: %w", err)
	}

	// Write to file
	if err := os.WriteFile(cs.filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write cookies: %w", err)
	}

	return nil
}

// Load loads cookies from file and sets them in the browser
func (cs *CookieStore) Load(page *rod.Page) error {
	// Check if file exists
	if _, err := os.Stat(cs.filePath); os.IsNotExist(err) {
		return nil // No cookies to load, not an error
	}

	// Read file
	data, err := os.ReadFile(cs.filePath)
	if err != nil {
		return fmt.Errorf("failed to read cookies: %w", err)
	}

	// Unmarshal JSON
	var storeCookies []Cookie
	if err := json.Unmarshal(data, &storeCookies); err != nil {
		return fmt.Errorf("failed to unmarshal cookies: %w", err)
	}

	// Convert to proto cookies and set in browser
	protoCookies := make([]*proto.NetworkCookieParam, 0, len(storeCookies))
	now := time.Now()

	for _, c := range storeCookies {
		// Skip expired cookies
		if !c.Expires.IsZero() && c.Expires.Before(now) {
			continue
		}

		sameSite := proto.NetworkCookieSameSiteNone
		switch c.SameSite {
		case "Strict":
			sameSite = proto.NetworkCookieSameSiteStrict
		case "Lax":
			sameSite = proto.NetworkCookieSameSiteLax
		}

		protoCookies = append(protoCookies, &proto.NetworkCookieParam{
			Name:     c.Name,
			Value:    c.Value,
			Domain:   c.Domain,
			Path:     c.Path,
			Expires:  proto.TimeSinceEpoch(c.Expires.Unix()),
			HTTPOnly: c.HTTPOnly,
			Secure:   c.Secure,
			SameSite: sameSite,
		})
	}

	// Set cookies in browser
	if err := page.SetCookies(protoCookies); err != nil {
		return fmt.Errorf("failed to set cookies: %w", err)
	}

	return nil
}

// Delete deletes the cookie file
func (cs *CookieStore) Delete() error {
	if _, err := os.Stat(cs.filePath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to delete
	}

	if err := os.Remove(cs.filePath); err != nil {
		return fmt.Errorf("failed to delete cookies: %w", err)
	}

	return nil
}

// Exists checks if cookie file exists
func (cs *CookieStore) Exists() bool {
	_, err := os.Stat(cs.filePath)
	return err == nil
}

package stealth

import (
	"github.com/go-rod/rod"
)

// Stealth manages all anti-detection techniques
type Stealth struct {
	config StealthConfig
}

// StealthConfig holds stealth configuration
type StealthConfig struct {
	DisableWebDriver    bool
	RandomizeUserAgent  bool
	RandomizeViewport   bool
	SpoofLocale         bool
	OverridePermissions bool
	RandomizeCanvas     bool
	RandomizeWebGL      bool
	UserAgent           string
	Timezone            string
	Language            string
}

// New creates a new Stealth manager
func New(cfg StealthConfig) *Stealth {
	return &Stealth{
		config: cfg,
	}
}

// Apply applies all enabled stealth techniques to the page
func (s *Stealth) Apply(page *rod.Page) error {
	if s.config.DisableWebDriver {
		if err := DisableWebDriver(page); err != nil {
			return err
		}
	}

	if s.config.SpoofLocale {
		if err := SpoofLocale(page, s.config.Timezone, s.config.Language); err != nil {
			return err
		}
	}

	if s.config.OverridePermissions {
		if err := OverridePermissions(page); err != nil {
			return err
		}
	}

	if s.config.RandomizeCanvas {
		if err := RandomizeCanvas(page); err != nil {
			return err
		}
	}

	if s.config.RandomizeWebGL {
		if err := RandomizeWebGL(page); err != nil {
			return err
		}
	}

	return nil
}

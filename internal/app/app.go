package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nkh/linkedin-automation/internal/auth"
	"github.com/nkh/linkedin-automation/internal/behavior"
	"github.com/nkh/linkedin-automation/internal/browser"
	"github.com/nkh/linkedin-automation/internal/config"
	"github.com/nkh/linkedin-automation/internal/scheduler"
	"github.com/nkh/linkedin-automation/internal/stealth"
	"github.com/nkh/linkedin-automation/internal/store"
)

// App represents the main application
type App struct {
	config    *config.Config
	logger    *Logger
	chrome    *browser.Chrome
	rod       *browser.Rod
	stealth   *stealth.Stealth
	behavior  *behavior.Engine
	auth      *auth.Authenticator
	store     store.Store
	scheduler *scheduler.Scheduler
	ctx       context.Context
	cancel    context.CancelFunc
}

// New creates a new application instance
func New(cfg *config.Config) (*App, error) {
	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	app := &App{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}

	return app, nil
}

// Initialize initializes all components
func (a *App) Initialize() error {
	var err error

	// Initialize logger
	a.logger, err = NewLogger(a.config.Paths.LogsDir)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	a.logger.Info("Initializing LinkedIn Automation Framework...")

	// Initialize store
	a.store, err = store.NewJSONStore(a.config.Paths.DataDir)
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	// Initialize scheduler
	a.scheduler = scheduler.New(a.config, a.store)

	// Launch Chrome
	a.logger.Info("Launching Chrome browser...")
	chromeConfig := browser.ChromeConfig{
		Headless:    a.config.Browser.Headless,
		Width:       a.config.Browser.Width,
		Height:      a.config.Browser.Height,
		UserDataDir: browser.GetUserDataDir(a.config.Paths.DataDir),
	}

	a.chrome, err = browser.LaunchChrome(a.ctx, chromeConfig)
	if err != nil {
		return fmt.Errorf("failed to launch Chrome: %w", err)
	}

	// Connect Rod
	a.logger.Info("Connecting to browser...")
	rodConfig := browser.RodConfig{
		ChromeURL: a.chrome.URL(),
		Timeout:   30 * time.Second,
	}

	a.rod, err = browser.NewRod(a.ctx, rodConfig)
	if err != nil {
		return fmt.Errorf("failed to connect Rod: %w", err)
	}

	// Create page
	page, err := a.rod.NewPage(a.ctx)
	if err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}

	// Apply stealth techniques
	a.logger.Info("Applying stealth techniques...")
	userAgent := stealth.GetUserAgent("", a.config.Stealth.RandomUserAgent)
	viewport := stealth.GetViewport(a.config.Browser.Width, a.config.Browser.Height, false)

	a.stealth = stealth.New(stealth.StealthConfig{
		DisableWebDriver:    true,
		RandomizeUserAgent:  a.config.Stealth.RandomUserAgent,
		RandomizeViewport:   false,
		SpoofLocale:         true,
		OverridePermissions: true,
		RandomizeCanvas:     true,
		RandomizeWebGL:      true,
		UserAgent:           userAgent,
		Timezone:            a.config.Stealth.Timezone,
		Language:            a.config.Stealth.Language,
	})

	if err := a.stealth.Apply(page); err != nil {
		return fmt.Errorf("failed to apply stealth: %w", err)
	}

	// Create behavior engine
	a.behavior = behavior.New(a.ctx, page)

	// Create authenticator
	authConfig := auth.Config{
		Email:      a.config.LinkedIn.Email,
		Password:   a.config.LinkedIn.Password,
		CookieFile: a.config.Paths.CookieFile,
	}

	a.auth = auth.New(authConfig, page, a.ctx)

	a.logger.Info("Initialization complete")
	return nil
}

// Run runs the main application loop
func (a *App) Run() error {
	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		a.logger.Info("Received shutdown signal")
		a.Shutdown()
	}()

	// Login
	a.logger.Info("Attempting login...")
	if err := a.auth.Login(a.ctx); err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	a.logger.Info("Successfully logged in to LinkedIn")

	// Main automation loop would go here
	// For now, just wait
	a.logger.Info("Bot is running. Press Ctrl+C to stop.")

	// Keep running until interrupted
	<-a.ctx.Done()

	return nil
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown() {
	a.logger.Info("Shutting down...")

	// Cancel context
	if a.cancel != nil {
		a.cancel()
	}

	// Close store
	if a.store != nil {
		if err := a.store.Close(); err != nil {
			a.logger.Error("Failed to close store: %v", err)
		}
	}

	// Close Rod
	if a.rod != nil {
		if err := a.rod.Close(); err != nil {
			a.logger.Error("Failed to close Rod: %v", err)
		}
	}

	// Close Chrome
	if a.chrome != nil {
		if err := a.chrome.Close(); err != nil {
			a.logger.Error("Failed to close Chrome: %v", err)
		}
	}

	// Close logger
	if a.logger != nil {
		a.logger.Info("Shutdown complete")
		if err := a.logger.Close(); err != nil {
			fmt.Printf("Failed to close logger: %v\n", err)
		}
	}
}

// GetLogger returns the logger
func (a *App) GetLogger() *Logger {
	return a.logger
}

// GetScheduler returns the scheduler
func (a *App) GetScheduler() *scheduler.Scheduler {
	return a.scheduler
}

// GetStore returns the store
func (a *App) GetStore() store.Store {
	return a.store
}

// GetBehavior returns the behavior engine
func (a *App) GetBehavior() *behavior.Engine {
	return a.behavior
}

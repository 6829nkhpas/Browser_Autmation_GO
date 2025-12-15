package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration
type Config struct {
	// LinkedIn credentials
	LinkedIn LinkedInConfig

	// Browser settings
	Browser BrowserConfig

	// Rate limiting
	RateLimit RateLimitConfig

	// Business hours
	BusinessHours BusinessHoursConfig

	// Stealth settings
	Stealth StealthConfig

	// Paths
	Paths PathsConfig

	// Search settings
	Search SearchConfig

	// Message settings
	Message MessageConfig

	// Break settings
	Breaks BreakConfig
}

// LinkedInConfig holds LinkedIn credentials
type LinkedInConfig struct {
	Email    string
	Password string
}

// BrowserConfig holds browser settings
type BrowserConfig struct {
	Headless bool
	Width    int
	Height   int
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	DailyConnectionLimit int
	HourlyMessageLimit   int
	MinActionDelayMs     int
	MaxActionDelayMs     int
}

// BusinessHoursConfig holds business hours configuration
type BusinessHoursConfig struct {
	Start int      // 24-hour format (0-23)
	End   int      // 24-hour format (0-23)
	Days  []string // e.g., ["Monday", "Tuesday", ...]
}

// StealthConfig holds stealth/anti-detection settings
type StealthConfig struct {
	Timezone         string
	Language         string
	RandomUserAgent  bool
}

// PathsConfig holds file paths
type PathsConfig struct {
	DataDir    string
	LogsDir    string
	CookieFile string
}

// SearchConfig holds search settings
type SearchConfig struct {
	MaxPages        int
	ResultsPerPage  int
}

// MessageConfig holds message settings
type MessageConfig struct {
	TemplateDir          string
	FollowUpDelayDaysMin int
	FollowUpDelayDaysMax int
}

// BreakConfig holds break scheduling settings
type BreakConfig struct {
	FrequencyMinutesMin int
	FrequencyMinutesMax int
	DurationMinutesMin  int
	DurationMinutesMax  int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		LinkedIn: LinkedInConfig{
			Email:    getEnv("LINKEDIN_EMAIL", ""),
			Password: getEnv("LINKEDIN_PASSWORD", ""),
		},
		Browser: BrowserConfig{
			Headless: getEnvBool("BROWSER_HEADLESS", false),
			Width:    getEnvInt("BROWSER_WIDTH", 1366),
			Height:   getEnvInt("BROWSER_HEIGHT", 768),
		},
		RateLimit: RateLimitConfig{
			DailyConnectionLimit: getEnvInt("DAILY_CONNECTION_LIMIT", 30),
			HourlyMessageLimit:   getEnvInt("HOURLY_MESSAGE_LIMIT", 10),
			MinActionDelayMs:     getEnvInt("MIN_ACTION_DELAY_MS", 2000),
			MaxActionDelayMs:     getEnvInt("MAX_ACTION_DELAY_MS", 5000),
		},
		BusinessHours: BusinessHoursConfig{
			Start: getEnvInt("BUSINESS_HOURS_START", 9),
			End:   getEnvInt("BUSINESS_HOURS_END", 17),
			Days:  getEnvStringSlice("BUSINESS_DAYS", []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}),
		},
		Stealth: StealthConfig{
			Timezone:        getEnv("TIMEZONE", time.Local.String()),
			Language:        getEnv("LANGUAGE", "en-US"),
			RandomUserAgent: getEnvBool("RANDOM_USER_AGENT", true),
		},
		Paths: PathsConfig{
			DataDir:    getEnv("DATA_DIR", "./data"),
			LogsDir:    getEnv("LOGS_DIR", "./logs"),
			CookieFile: getEnv("COOKIE_FILE", "./data/cookies.json"),
		},
		Search: SearchConfig{
			MaxPages:       getEnvInt("SEARCH_MAX_PAGES", 5),
			ResultsPerPage: getEnvInt("SEARCH_RESULTS_PER_PAGE", 10),
		},
		Message: MessageConfig{
			TemplateDir:          getEnv("MESSAGE_TEMPLATE_DIR", "./assets/templates"),
			FollowUpDelayDaysMin: getEnvInt("FOLLOW_UP_DELAY_DAYS_MIN", 1),
			FollowUpDelayDaysMax: getEnvInt("FOLLOW_UP_DELAY_DAYS_MAX", 3),
		},
		Breaks: BreakConfig{
			FrequencyMinutesMin: getEnvInt("BREAK_FREQUENCY_MINUTES_MIN", 30),
			FrequencyMinutesMax: getEnvInt("BREAK_FREQUENCY_MINUTES_MAX", 60),
			DurationMinutesMin:  getEnvInt("BREAK_DURATION_MINUTES_MIN", 5),
			DurationMinutesMax:  getEnvInt("BREAK_DURATION_MINUTES_MAX", 20),
		},
	}

	if err := Validate(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvStringSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
}

package config

import (
	"fmt"
)

// Validate validates the configuration
func Validate(cfg *Config) error {
	// Validate LinkedIn credentials
	if cfg.LinkedIn.Email == "" {
		return fmt.Errorf("LINKEDIN_EMAIL is required")
	}
	if cfg.LinkedIn.Password == "" {
		return fmt.Errorf("LINKEDIN_PASSWORD is required")
	}

	// Validate browser settings
	if cfg.Browser.Width < 800 || cfg.Browser.Width > 3840 {
		return fmt.Errorf("BROWSER_WIDTH must be between 800 and 3840, got %d", cfg.Browser.Width)
	}
	if cfg.Browser.Height < 600 || cfg.Browser.Height > 2160 {
		return fmt.Errorf("BROWSER_HEIGHT must be between 600 and 2160, got %d", cfg.Browser.Height)
	}

	// Validate rate limits
	if cfg.RateLimit.DailyConnectionLimit < 1 || cfg.RateLimit.DailyConnectionLimit > 100 {
		return fmt.Errorf("DAILY_CONNECTION_LIMIT must be between 1 and 100, got %d", cfg.RateLimit.DailyConnectionLimit)
	}
	if cfg.RateLimit.HourlyMessageLimit < 1 || cfg.RateLimit.HourlyMessageLimit > 50 {
		return fmt.Errorf("HOURLY_MESSAGE_LIMIT must be between 1 and 50, got %d", cfg.RateLimit.HourlyMessageLimit)
	}
	if cfg.RateLimit.MinActionDelayMs < 500 || cfg.RateLimit.MinActionDelayMs > 10000 {
		return fmt.Errorf("MIN_ACTION_DELAY_MS must be between 500 and 10000, got %d", cfg.RateLimit.MinActionDelayMs)
	}
	if cfg.RateLimit.MaxActionDelayMs < cfg.RateLimit.MinActionDelayMs {
		return fmt.Errorf("MAX_ACTION_DELAY_MS (%d) must be >= MIN_ACTION_DELAY_MS (%d)", 
			cfg.RateLimit.MaxActionDelayMs, cfg.RateLimit.MinActionDelayMs)
	}

	// Validate business hours
	if cfg.BusinessHours.Start < 0 || cfg.BusinessHours.Start > 23 {
		return fmt.Errorf("BUSINESS_HOURS_START must be between 0 and 23, got %d", cfg.BusinessHours.Start)
	}
	if cfg.BusinessHours.End < 0 || cfg.BusinessHours.End > 23 {
		return fmt.Errorf("BUSINESS_HOURS_END must be between 0 and 23, got %d", cfg.BusinessHours.End)
	}
	if cfg.BusinessHours.Start >= cfg.BusinessHours.End {
		return fmt.Errorf("BUSINESS_HOURS_START (%d) must be < BUSINESS_HOURS_END (%d)", 
			cfg.BusinessHours.Start, cfg.BusinessHours.End)
	}
	if len(cfg.BusinessHours.Days) == 0 {
		return fmt.Errorf("BUSINESS_DAYS must contain at least one day")
	}

	// Validate valid day names
	validDays := map[string]bool{
		"Monday": true, "Tuesday": true, "Wednesday": true, "Thursday": true,
		"Friday": true, "Saturday": true, "Sunday": true,
	}
	for _, day := range cfg.BusinessHours.Days {
		if !validDays[day] {
			return fmt.Errorf("invalid day in BUSINESS_DAYS: %s", day)
		}
	}

	// Validate search settings
	if cfg.Search.MaxPages < 1 || cfg.Search.MaxPages > 50 {
		return fmt.Errorf("SEARCH_MAX_PAGES must be between 1 and 50, got %d", cfg.Search.MaxPages)
	}

	// Validate message settings
	if cfg.Message.FollowUpDelayDaysMin < 1 || cfg.Message.FollowUpDelayDaysMin > 30 {
		return fmt.Errorf("FOLLOW_UP_DELAY_DAYS_MIN must be between 1 and 30, got %d", cfg.Message.FollowUpDelayDaysMin)
	}
	if cfg.Message.FollowUpDelayDaysMax < cfg.Message.FollowUpDelayDaysMin {
		return fmt.Errorf("FOLLOW_UP_DELAY_DAYS_MAX (%d) must be >= FOLLOW_UP_DELAY_DAYS_MIN (%d)", 
			cfg.Message.FollowUpDelayDaysMax, cfg.Message.FollowUpDelayDaysMin)
	}

	// Validate break settings
	if cfg.Breaks.FrequencyMinutesMin < 10 || cfg.Breaks.FrequencyMinutesMin > 240 {
		return fmt.Errorf("BREAK_FREQUENCY_MINUTES_MIN must be between 10 and 240, got %d", cfg.Breaks.FrequencyMinutesMin)
	}
	if cfg.Breaks.FrequencyMinutesMax < cfg.Breaks.FrequencyMinutesMin {
		return fmt.Errorf("BREAK_FREQUENCY_MINUTES_MAX (%d) must be >= BREAK_FREQUENCY_MINUTES_MIN (%d)", 
			cfg.Breaks.FrequencyMinutesMax, cfg.Breaks.FrequencyMinutesMin)
	}
	if cfg.Breaks.DurationMinutesMin < 1 || cfg.Breaks.DurationMinutesMin > 60 {
		return fmt.Errorf("BREAK_DURATION_MINUTES_MIN must be between 1 and 60, got %d", cfg.Breaks.DurationMinutesMin)
	}
	if cfg.Breaks.DurationMinutesMax < cfg.Breaks.DurationMinutesMin {
		return fmt.Errorf("BREAK_DURATION_MINUTES_MAX (%d) must be >= BREAK_DURATION_MINUTES_MIN (%d)", 
			cfg.Breaks.DurationMinutesMax, cfg.Breaks.DurationMinutesMin)
	}

	return nil
}

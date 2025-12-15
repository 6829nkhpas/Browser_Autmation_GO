package scheduler

import (
	"fmt"
	"time"

	"github.com/nkh/linkedin-automation/internal/behavior"
	"github.com/nkh/linkedin-automation/internal/config"
	"github.com/nkh/linkedin-automation/internal/store"
)

// Scheduler manages rate limiting and business hours
type Scheduler struct {
	config *config.Config
	store  store.Store
}

// New creates a new scheduler
func New(cfg *config.Config, st store.Store) *Scheduler {
	return &Scheduler{
		config: cfg,
		store:  st,
	}
}

// CanConnect checks if a connection request can be made
func (s *Scheduler) CanConnect() (bool, string) {
	// Check business hours
	if !s.IsBusinessHours() {
		return false, "outside business hours"
	}

	// Check daily limit
	count, err := s.store.GetDailyActionCount(store.ActionConnectionRequest)
	if err != nil {
		return false, fmt.Sprintf("failed to get action count: %v", err)
	}

	if count >= s.config.RateLimit.DailyConnectionLimit {
		return false, fmt.Sprintf("daily connection limit reached (%d/%d)", count, s.config.RateLimit.DailyConnectionLimit)
	}

	return true, ""
}

// CanMessage checks if a message can be sent
func (s *Scheduler) CanMessage() (bool, string) {
	// Check business hours
	if !s.IsBusinessHours() {
		return false, "outside business hours"
	}

	// Check hourly limit
	count, err := s.store.GetHourlyActionCount(store.ActionMessageSent)
	if err != nil {
		return false, fmt.Sprintf("failed to get action count: %v", err)
	}

	if count >= s.config.RateLimit.HourlyMessageLimit {
		return false, fmt.Sprintf("hourly message limit reached (%d/%d)", count, s.config.RateLimit.HourlyMessageLimit)
	}

	return true, ""
}

// IsBusinessHours checks if current time is within business hours
func (s *Scheduler) IsBusinessHours() bool {
	now := time.Now()

	// Check day of week
	dayName := now.Weekday().String()
	isBusinessDay := false
	for _, day := range s.config.BusinessHours.Days {
		if day == dayName {
			isBusinessDay = true
			break
		}
	}

	if !isBusinessDay {
		return false
	}

	// Check hour
	hour := now.Hour()
	if hour < s.config.BusinessHours.Start || hour >= s.config.BusinessHours.End {
		return false
	}

	return true
}

// WaitForBusinessHours waits until business hours
func (s *Scheduler) WaitForBusinessHours() {
	for !s.IsBusinessHours() {
		// Wait 15 minutes and check again
		time.Sleep(15 * time.Minute)
	}
}

// ApplyActionDelay applies the configured delay between actions
func (s *Scheduler) ApplyActionDelay() {
	behavior.WaitHuman(
		s.config.RateLimit.MinActionDelayMs,
		s.config.RateLimit.MaxActionDelayMs,
	)
}

// ShouldTakeBreak determines if a break should be taken
func (s *Scheduler) ShouldTakeBreak(actionsSinceBreak int) bool {
	// Take break every 10-20 actions
	threshold := behavior.GetRandomInRange(10, 20)
	return actionsSinceBreak >= threshold
}

// TakeBreak pauses execution for a random break duration
func (s *Scheduler) TakeBreak() {
	durationMinutes := behavior.GetRandomInRange(
		s.config.Breaks.DurationMinutesMin,
		s.config.Breaks.DurationMinutesMax,
	)

	fmt.Printf("Taking break for %d minutes...\n", durationMinutes)
	time.Sleep(time.Duration(durationMinutes) * time.Minute)
}

// GetNextBreakTime calculates when the next break should occur
func (s *Scheduler) GetNextBreakTime() time.Time {
	minutesUntilBreak := behavior.GetRandomInRange(
		s.config.Breaks.FrequencyMinutesMin,
		s.config.Breaks.FrequencyMinutesMax,
	)

	return time.Now().Add(time.Duration(minutesUntilBreak) * time.Minute)
}

// GetDailyStats returns statistics for today
func (s *Scheduler) GetDailyStats() (map[store.ActionType]int, error) {
	stats := make(map[store.ActionType]int)

	actionTypes := []store.ActionType{
		store.ActionProfileVisit,
		store.ActionConnectionRequest,
		store.ActionMessageSent,
		store.ActionSearch,
	}

	for _, actionType := range actionTypes {
		count, err := s.store.GetDailyActionCount(actionType)
		if err != nil {
			return nil, err
		}
		stats[actionType] = count
	}

	return stats, nil
}

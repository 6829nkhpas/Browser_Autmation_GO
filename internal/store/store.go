package store

import (
	"time"
)

// ActionType represents the type of action performed
type ActionType string

const (
	ActionProfileVisit      ActionType = "profile_visit"
	ActionConnectionRequest ActionType = "connection_request"
	ActionMessageSent       ActionType = "message_sent"
	ActionSearch            ActionType = "search"
)

// Action represents a recorded action
type Action struct {
	ID          string     `json:"id"`
	Type        ActionType `json:"type"`
	ProfileURL  string     `json:"profile_url,omitempty"`
	ProfileName string     `json:"profile_name,omitempty"`
	Message     string     `json:"message,omitempty"`
	Timestamp   time.Time  `json:"timestamp"`
	Success     bool       `json:"success"`
	ErrorMsg    string     `json:"error_msg,omitempty"`
}

// Profile represents profile data
type Profile struct {
	URL         string    `json:"url"`
	Name        string    `json:"name"`
	Headline    string    `json:"headline"`
	Location    string    `json:"location"`
	Company     string    `json:"company,omitempty"`
	FirstSeen   time.Time `json:"first_seen"`
	LastVisited time.Time `json:"last_visited"`
	Connected   bool      `json:"connected"`
	MessageSent bool      `json:"message_sent"`
}

// Store defines the interface for persistence
type Store interface {
	// Actions
	SaveAction(action Action) error
	GetActions(actionType ActionType, limit int) ([]Action, error)
	GetActionsByDate(date time.Time) ([]Action, error)

	// Profiles
	SaveProfile(profile Profile) error
	GetProfile(url string) (*Profile, error)
	ProfileExists(url string) bool
	UpdateProfile(profile Profile) error

	// Statistics
	GetActionCount(actionType ActionType, since time.Time) (int, error)
	GetDailyActionCount(actionType ActionType) (int, error)
	GetHourlyActionCount(actionType ActionType) (int, error)

	// Cleanup
	Close() error
}

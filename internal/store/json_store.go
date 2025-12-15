package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// JSONStore implements Store using JSON files
type JSONStore struct {
	dataDir      string
	actions      []Action
	profiles     map[string]Profile
	mutex        sync.RWMutex
	actionsFile  string
	profilesFile string
}

// NewJSONStore creates a new JSON-based store
func NewJSONStore(dataDir string) (*JSONStore, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	store := &JSONStore{
		dataDir:      dataDir,
		profiles:     make(map[string]Profile),
		actionsFile:  filepath.Join(dataDir, "actions.json"),
		profilesFile: filepath.Join(dataDir, "profiles.json"),
	}

	// Load existing data
	if err := store.load(); err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	return store, nil
}

// SaveAction saves an action
func (s *JSONStore) SaveAction(action Action) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Generate ID if not set
	if action.ID == "" {
		action.ID = fmt.Sprintf("%d-%s", time.Now().Unix(), action.Type)
	}

	s.actions = append(s.actions, action)

	return s.saveActions()
}

// GetActions retrieves actions by type
func (s *JSONStore) GetActions(actionType ActionType, limit int) ([]Action, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var filtered []Action
	for i := len(s.actions) - 1; i >= 0 && len(filtered) < limit; i-- {
		if s.actions[i].Type == actionType {
			filtered = append(filtered, s.actions[i])
		}
	}

	return filtered, nil
}

// GetActionsByDate retrieves actions for a specific date
func (s *JSONStore) GetActionsByDate(date time.Time) ([]Action, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var filtered []Action
	for _, action := range s.actions {
		if action.Timestamp.After(startOfDay) && action.Timestamp.Before(endOfDay) {
			filtered = append(filtered, action)
		}
	}

	return filtered, nil
}

// SaveProfile saves a profile
func (s *JSONStore) SaveProfile(profile Profile) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.profiles[profile.URL] = profile

	return s.saveProfiles()
}

// GetProfile retrieves a profile by URL
func (s *JSONStore) GetProfile(url string) (*Profile, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	profile, exists := s.profiles[url]
	if !exists {
		return nil, fmt.Errorf("profile not found")
	}

	return &profile, nil
}

// ProfileExists checks if a profile exists
func (s *JSONStore) ProfileExists(url string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, exists := s.profiles[url]
	return exists
}

// UpdateProfile updates an existing profile
func (s *JSONStore) UpdateProfile(profile Profile) error {
	return s.SaveProfile(profile)
}

// GetActionCount returns count of actions since a specific time
func (s *JSONStore) GetActionCount(actionType ActionType, since time.Time) (int, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	count := 0
	for _, action := range s.actions {
		if action.Type == actionType && action.Timestamp.After(since) && action.Success {
			count++
		}
	}

	return count, nil
}

// GetDailyActionCount returns count of actions today
func (s *JSONStore) GetDailyActionCount(actionType ActionType) (int, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return s.GetActionCount(actionType, startOfDay)
}

// GetHourlyActionCount returns count of actions in the last hour
func (s *JSONStore) GetHourlyActionCount(actionType ActionType) (int, error) {
	return s.GetActionCount(actionType, time.Now().Add(-1*time.Hour))
}

// Close closes the store
func (s *JSONStore) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Save any pending data
	if err := s.saveActions(); err != nil {
		return err
	}
	if err := s.saveProfiles(); err != nil {
		return err
	}

	return nil
}

// load loads data from files
func (s *JSONStore) load() error {
	// Load actions
	if _, err := os.Stat(s.actionsFile); err == nil {
		data, err := os.ReadFile(s.actionsFile)
		if err != nil {
			return fmt.Errorf("failed to read actions file: %w", err)
		}

		if err := json.Unmarshal(data, &s.actions); err != nil {
			return fmt.Errorf("failed to unmarshal actions: %w", err)
		}
	}

	// Load profiles
	if _, err := os.Stat(s.profilesFile); err == nil {
		data, err := os.ReadFile(s.profilesFile)
		if err != nil {
			return fmt.Errorf("failed to read profiles file: %w", err)
		}

		if err := json.Unmarshal(data, &s.profiles); err != nil {
			return fmt.Errorf("failed to unmarshal profiles: %w", err)
		}
	}

	return nil
}

// saveActions saves actions to file
func (s *JSONStore) saveActions() error {
	data, err := json.MarshalIndent(s.actions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal actions: %w", err)
	}

	if err := os.WriteFile(s.actionsFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write actions file: %w", err)
	}

	return nil
}

// saveProfiles saves profiles to file
func (s *JSONStore) saveProfiles() error {
	data, err := json.MarshalIndent(s.profiles, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profiles: %w", err)
	}

	if err := os.WriteFile(s.profilesFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write profiles file: %w", err)
	}

	return nil
}

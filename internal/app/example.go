package app

import (
	"fmt"
	"time"

	"github.com/nkh/linkedin-automation/internal/connect"
	"github.com/nkh/linkedin-automation/internal/message"
	"github.com/nkh/linkedin-automation/internal/profile"
	"github.com/nkh/linkedin-automation/internal/search"
	"github.com/nkh/linkedin-automation/internal/store"
)

// runExampleAutomation demonstrates how to use all components together
// This is an example flow - customize it for your specific needs
func (a *App) runExampleAutomation() error {
	a.logger.Info("Starting example automation flow...")

	// Initialize components
	searchEngine := search.New(a.ctx, a.behavior, a.store)
	profileVisitor := profile.New(a.ctx, a.behavior, a.store)
	connector := connect.New(a.ctx, a.behavior, profileVisitor, a.store)
	messenger, err := message.New(a.ctx, a.behavior, a.store, a.config.Message.TemplateDir)
	if err != nil {
		return fmt.Errorf("failed to create messenger: %w", err)
	}

	// Example 1: Search for people
	a.logger.Info("Searching for profiles...")
	profiles, err := searchEngine.Search(search.Config{
		Keywords: "Software Engineer",
		JobTitle: "Senior Engineer",
		Location: "San Francisco",
		MaxPages: 2, // Search first 2 pages
	})
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	a.logger.Info("Found %d profiles", len(profiles))

	// Example 2: Visit and connect with profiles
	connectionsToday := 0
	maxConnectionsPerDay := a.config.RateLimit.DailyConnectionLimit

	for _, profileURL := range profiles {
		// Check if we've hit daily limit
		if connectionsToday >= maxConnectionsPerDay {
			a.logger.Info("Daily connection limit reached (%d/%d)", connectionsToday, maxConnectionsPerDay)
			break
		}

		// Check scheduler permissions
		canConnect, reason := a.scheduler.CanConnect()
		if !canConnect {
			a.logger.RateLimit("CONNECTION", reason)
			break
		}

		// Check if already connected
		if connector.IsAlreadyConnected(profileURL) {
			a.logger.Info("Already connected to profile, skipping...")
			continue
		}

		// Extract profile data for personalization
		profileData, err := profileVisitor.Extract(profileURL)
		if err != nil {
			a.logger.Error("Failed to extract profile data: %v", err)
			continue
		}

		// Generate personalized note
		variables := map[string]string{
			"name":    profileData.Name,
			"company": profileData.Location, // Could parse company from headline
			"field":   "software engineering",
		}
		note := messenger.GetConnectionRequestTemplate(variables)

		// Send connection request
		a.logger.Info("Sending connection request to: %s", profileData.Name)
		if err := connector.Connect(profileURL, note); err != nil {
			a.logger.Error("Connection request failed: %v", err)
			continue
		}

		connectionsToday++
		a.logger.Action("CONNECTION_REQUEST", profileURL, true)

		// Apply rate limit delay
		a.scheduler.ApplyActionDelay()

		// Take break after every 5-10 actions
		if a.scheduler.ShouldTakeBreak(connectionsToday) {
			a.logger.Info("Taking a natural break...")
			a.scheduler.TakeBreak()
		}
	}

	// Example 3: Follow up with accepted connections
	a.logger.Info("Checking for accepted connections to follow up...")

	// Get recent connection requests
	recentActions, err := a.store.GetActions(store.ActionConnectionRequest, 50)
	if err != nil {
		a.logger.Error("Failed to get recent actions: %v", err)
	} else {
		for _, action := range recentActions {
			// Check if enough time has passed (1-3 days)
			daysSince := time.Since(action.Timestamp).Hours() / 24
			minDays := float64(a.config.Message.FollowUpDelayDaysMin)
			maxDays := float64(a.config.Message.FollowUpDelayDaysMax)

			if daysSince >= minDays && daysSince <= maxDays {
				// Check if we can send message
				canMessage, reason := a.scheduler.CanMessage()
				if !canMessage {
					a.logger.RateLimit("MESSAGE", reason)
					break
				}

				// Check if already messaged
				profile, _ := a.store.GetProfile(action.ProfileURL)
				if profile != nil && profile.MessageSent {
					continue
				}

				// Send follow-up message
				variables := map[string]string{
					"name": action.ProfileName,
				}

				a.logger.Info("Sending follow-up message to: %s", action.ProfileName)
				if err := messenger.SendFollowUp(action.ProfileURL, variables); err != nil {
					a.logger.Error("Failed to send message: %v", err)
					continue
				}

				// Update profile
				if profile != nil {
					profile.MessageSent = true
					_ = a.store.UpdateProfile(*profile)
				}

				a.logger.Action("MESSAGE_SENT", action.ProfileURL, true)
				a.scheduler.ApplyActionDelay()
			}
		}
	}

	// Print daily statistics
	stats, err := a.scheduler.GetDailyStats()
	if err != nil {
		a.logger.Error("Failed to get stats: %v", err)
	} else {
		a.logger.Info("Daily Statistics:")
		a.logger.Info("  Profile Visits: %d", stats[store.ActionProfileVisit])
		a.logger.Info("  Connection Requests: %d/%d", stats[store.ActionConnectionRequest], maxConnectionsPerDay)
		a.logger.Info("  Messages Sent: %d", stats[store.ActionMessageSent])
	}

	a.logger.Info("Example automation flow complete!")
	return nil
}

package search

import (
	"context"
	"fmt"
	"time"

	"github.com/nkh/linkedin-automation/internal/behavior"
	"github.com/nkh/linkedin-automation/internal/linkedin"
	"github.com/nkh/linkedin-automation/internal/store"
)

// Engine handles LinkedIn people search
type Engine struct {
	behavior *behavior.Engine
	store    store.Store
	ctx      context.Context
}

// Config holds search configuration
type Config struct {
	Keywords string
	JobTitle string
	Company  string
	Location string
	MaxPages int
}

// New creates a new search engine
func New(ctx context.Context, behaviorEng *behavior.Engine, st store.Store) *Engine {
	return &Engine{
		behavior: behaviorEng,
		store:    st,
		ctx:      ctx,
	}
}

// Search performs a people search and returns profile URLs
func (e *Engine) Search(cfg Config) ([]string, error) {
	// Build search URL
	url := buildSearchURL(cfg)

	// Navigate to search page
	if err := e.behavior.Navigate(url); err != nil {
		return nil, fmt.Errorf("failed to navigate to search: %w", err)
	}

	// Wait for results to load
	if err := e.behavior.WaitForElement(linkedin.SearchResultsList, 10*time.Second); err != nil {
		return nil, fmt.Errorf("search results didn't load: %w", err)
	}

	var allProfiles []string
	seenProfiles := make(map[string]bool)

	// Paginate through results
	for page := 1; page <= cfg.MaxPages; page++ {
		// Extract profiles from current page
		profiles, err := e.extractProfilesFromPage()
		if err != nil {
			return nil, fmt.Errorf("failed to extract profiles from page %d: %w", page, err)
		}

		// Deduplicate
		for _, profile := range profiles {
			if !seenProfiles[profile] {
				seenProfiles[profile] = true
				allProfiles = append(allProfiles, profile)
			}
		}

		// Natural scrolling to simulate reading results
		e.behavior.Scroll(behavior.GetRandomInRange(300, 600))
		behavior.WaitHuman(2000, 4000)

		// Check for next page
		if page < cfg.MaxPages {
			hasNext, err := e.goToNextPage()
			if err != nil {
				return nil, fmt.Errorf("failed to navigate to next page: %w", err)
			}
			if !hasNext {
				break // No more pages
			}
		}

		// Record search action
		_ = e.store.SaveAction(store.Action{
			Type:      store.ActionSearch,
			Timestamp: time.Now(),
			Success:   true,
		})
	}

	return allProfiles, nil
}

// extractProfilesFromPage extracts profile URLs from current search results page
func (e *Engine) extractProfilesFromPage() ([]string, error) {
	page := e.behavior.Page()

	// Find all result items
	elements, err := page.Elements(linkedin.SearchResultItem)
	if err != nil {
		return nil, fmt.Errorf("failed to find result items: %w", err)
	}

	var profiles []string
	for _, elem := range elements {
		// Find profile link within result item
		link, err := elem.Element(linkedin.SearchResultLink)
		if err != nil {
			continue // Skip if no link found
		}

		// Get href attribute
		href, err := link.Property("href")
		if err != nil {
			continue
		}

		profileURL := href.String()
		if profileURL != "" && contains(profileURL, "/in/") {
			profiles = append(profiles, profileURL)
		}
	}

	return profiles, nil
}

// goToNextPage navigates to the next page of search results
func (e *Engine) goToNextPage() (bool, error) {
	// Check if next button exists
	if !e.behavior.HasElement(linkedin.NextPageButton) {
		return false, nil // No more pages
	}

	// Click next button
	if err := e.behavior.Click(linkedin.NextPageButton); err != nil {
		return false, fmt.Errorf("failed to click next button: %w", err)
	}

	// Wait for new results to load
	behavior.WaitHuman(2000, 4000)

	return true, nil
}

// buildSearchURL builds the LinkedIn search URL with filters
func buildSearchURL(cfg Config) string {
	filters := make(map[string]string)

	if cfg.JobTitle != "" {
		filters["title"] = cfg.JobTitle
	}
	if cfg.Company != "" {
		filters["company"] = cfg.Company
	}
	if cfg.Location != "" {
		filters["location"] = cfg.Location
	}

	return linkedin.BuildSearchURL(cfg.Keywords, filters)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

package profile

import (
	"context"
	"fmt"
	"time"

	"github.com/nkh/linkedin-automation/internal/behavior"
	"github.com/nkh/linkedin-automation/internal/linkedin"
	"github.com/nkh/linkedin-automation/internal/store"
)

// Visitor handles profile visiting with natural behavior
type Visitor struct {
	behavior *behavior.Engine
	store    store.Store
	ctx      context.Context
}

// New creates a new profile visitor
func New(ctx context.Context, behaviorEng *behavior.Engine, st store.Store) *Visitor {
	return &Visitor{
		behavior: behaviorEng,
		store:    st,
		ctx:      ctx,
	}
}

// Visit visits a profile with natural human-like behavior
func (v *Visitor) Visit(profileURL string) error {
	// Navigate to profile
	if err := v.behavior.Navigate(profileURL); err != nil {
		return fmt.Errorf("failed to navigate to profile: %w", err)
	}

	// Wait for profile to load
	if err := v.behavior.WaitForElement(linkedin.ProfileName, 10*time.Second); err != nil {
		return fmt.Errorf("profile didn't load: %w", err)
	}

	// Simulate reading the profile name/headline
	behavior.WaitReading(10) // ~3 seconds

	// Natural scrolling through profile
	if err := v.scrollProfile(); err != nil {
		return fmt.Errorf("failed to scroll profile: %w", err)
	}

	// Hover over sections (simulate interest)
	v.hoverSections()

	// Thinking pause before any action
	behavior.WaitThinking()

	// Record visit
	_ = v.store.SaveAction(store.Action{
		Type:       store.ActionProfileVisit,
		ProfileURL: profileURL,
		Timestamp:  time.Now(),
		Success:    true,
	})

	return nil
}

// scrollProfile scrolls through the profile naturally
func (v *Visitor) scrollProfile() error {
	// Get approximate page height
	page := v.behavior.Page()
	height, err := page.Eval(`() => document.body.scrollHeight`)
	if err != nil {
		return fmt.Errorf("failed to get page height: %w", err)
	}

	totalHeight := int(height.Value.Num())
	scrolled := 0

	// Scroll in chunks with reading pauses
	for scrolled < totalHeight && scrolled < 3000 { // Max 3000px scroll
		chunkSize := behavior.GetRandomInRange(300, 600)
		if err := v.behavior.Scroll(chunkSize); err != nil {
			return err
		}

		scrolled += chunkSize

		// Reading pause (1-3 seconds)
		behavior.WaitHuman(1000, 3000)

		// Occasionally scroll back up (re-reading)
		if behavior.ShouldTakeBreak(15) { // 15% chance
			backScroll := -behavior.GetRandomInRange(50, 150)
			v.behavior.Scroll(backScroll)
			behavior.WaitHuman(500, 1000)
		}
	}

	return nil
}

// hoverSections hovers over profile sections to simulate interest
func (v *Visitor) hoverSections() {
	sections := []string{
		linkedin.ExperienceSection,
		linkedin.EducationSection,
		linkedin.SkillsSection,
	}

	// Hover over 1-2 random sections
	numSections := behavior.GetRandomInRange(1, 2)
	for i := 0; i < numSections && i < len(sections); i++ {
		section := sections[behavior.GetRandomInRange(0, len(sections)-1)]

		if v.behavior.HasElement(section) {
			_ = v.behavior.Hover(section)
			behavior.WaitHuman(500, 1500)
		}
	}
}

// Extract extracts profile data
func (v *Visitor) Extract(profileURL string) (*store.Profile, error) {
	page := v.behavior.Page()

	profile := &store.Profile{
		URL:       profileURL,
		FirstSeen: time.Now(),
	}

	// Extract name
	if nameElem, err := page.Element(linkedin.ProfileName); err == nil {
		if text, err := nameElem.Text(); err == nil {
			profile.Name = text
		}
	}

	// Extract headline
	if headlineElem, err := page.Element(linkedin.ProfileHeadline); err == nil {
		if text, err := headlineElem.Text(); err == nil {
			profile.Headline = text
		}
	}

	// Extract location
	if locationElem, err := page.Element(linkedin.ProfileLocation); err == nil {
		if text, err := locationElem.Text(); err == nil {
			profile.Location = text
		}
	}

	profile.LastVisited = time.Now()

	return profile, nil
}

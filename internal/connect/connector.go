package connect

import (
	"context"
	"fmt"
	"time"

	"github.com/nkh/linkedin-automation/internal/behavior"
	"github.com/nkh/linkedin-automation/internal/linkedin"
	"github.com/nkh/linkedin-automation/internal/profile"
	"github.com/nkh/linkedin-automation/internal/store"
)

// Connector handles connection requests
type Connector struct {
	behavior *behavior.Engine
	visitor  *profile.Visitor
	store    store.Store
	ctx      context.Context
}

// New creates a new connector
func New(ctx context.Context, behaviorEng *behavior.Engine, vis *profile.Visitor, st store.Store) *Connector {
	return &Connector{
		behavior: behaviorEng,
		visitor:  vis,
		store:    st,
		ctx:      ctx,
	}
}

// Connect sends a connection request to a profile
// Flow: Visit → Scroll → Hover → Pause → Connect → Note → Send
func (c *Connector) Connect(profileURL string, note string) error {
	// Step 1: Visit profile naturally
	if err := c.visitor.Visit(profileURL); err != nil {
		return fmt.Errorf("failed to visit profile: %w", err)
	}

	// Step 2: Thinking pause (deciding to connect)
	behavior.WaitThinking()

	// Step 3: Check if connect button exists
	if !c.behavior.HasElement(linkedin.ConnectButton) && !c.behavior.HasElement(linkedin.ConnectButtonAlt) {
		return fmt.Errorf("connect button not found - may already be connected")
	}

	// Step 4: Hover over connect button
	selector := linkedin.ConnectButton
	if !c.behavior.HasElement(selector) {
		selector = linkedin.ConnectButtonAlt
	}

	if err := c.behavior.Hover(selector); err != nil {
		return fmt.Errorf("failed to hover connect button: %w", err)
	}

	behavior.WaitHuman(300, 700)

	// Step 5: Click connect
	if err := c.behavior.Click(selector); err != nil {
		return fmt.Errorf("failed to click connect: %w", err)
	}

	// Wait for modal to appear
	behavior.WaitHuman(1000, 2000)

	// Step 6: Add note if provided
	if note != "" {
		if err := c.addNote(note); err != nil {
			return fmt.Errorf("failed to add note: %w", err)
		}
	}

	// Step 7: Send invitation
	if err := c.sendInvitation(); err != nil {
		return fmt.Errorf("failed to send invitation: %w", err)
	}

	// Step 8: Record action
	profileData, _ := c.visitor.Extract(profileURL)
	profileName := ""
	if profileData != nil {
		profileName = profileData.Name
	}

	_ = c.store.SaveAction(store.Action{
		Type:        store.ActionConnectionRequest,
		ProfileURL:  profileURL,
		ProfileName: profileName,
		Message:     note,
		Timestamp:   time.Now(),
		Success:     true,
	})

	// Step 9: Post-action cooldown
	behavior.WaitHuman(2000, 4000)

	return nil
}

// addNote adds a personalized note to the connection request
func (c *Connector) addNote(note string) error {
	// Check if "Add note" button exists
	if c.behavior.HasElement(linkedin.AddNoteButton) {
		if err := c.behavior.Click(linkedin.AddNoteButton); err != nil {
			return fmt.Errorf("failed to click add note: %w", err)
		}

		behavior.WaitHuman(500, 1000)
	}

	// Type note in textarea
	if !c.behavior.HasElement(linkedin.NoteTextarea) {
		return fmt.Errorf("note textarea not found")
	}

	// Truncate note to 300 characters (LinkedIn limit)
	if len(note) > 300 {
		note = note[:297] + "..."
	}

	if err := c.behavior.Type(linkedin.NoteTextarea, note); err != nil {
		return fmt.Errorf("failed to type note: %w", err)
	}

	behavior.WaitHuman(500, 1000)

	return nil
}

// sendInvitation sends the connection invitation
func (c *Connector) sendInvitation() error {
	// Try to find send button
	sendSelectors := []string{
		linkedin.SendInviteButton,
		linkedin.SendButton,
	}

	var sendSelector string
	for _, selector := range sendSelectors {
		if c.behavior.HasElement(selector) {
			sendSelector = selector
			break
		}
	}

	if sendSelector == "" {
		return fmt.Errorf("send button not found")
	}

	// Click send
	if err := c.behavior.Click(sendSelector); err != nil {
		return fmt.Errorf("failed to click send: %w", err)
	}

	// Wait for confirmation
	behavior.WaitHuman(1500, 2500)

	return nil
}

// IsAlreadyConnected checks if already connected to a profile
func (c *Connector) IsAlreadyConnected(profileURL string) bool {
	// Navigate to profile
	if err := c.behavior.Navigate(profileURL); err != nil {
		return false
	}

	// Check for message button (indicates connection)
	if c.behavior.HasElement(linkedin.MessageButton) {
		return true
	}

	// Check for pending button
	if c.behavior.HasElement(linkedin.PendingButton) {
		return true
	}

	return false
}

package message

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nkh/linkedin-automation/internal/behavior"
	"github.com/nkh/linkedin-automation/internal/linkedin"
	"github.com/nkh/linkedin-automation/internal/store"
)

// Messenger handles LinkedIn messaging
type Messenger struct {
	behavior  *behavior.Engine
	store     store.Store
	templates *Templates
	ctx       context.Context
}

// New creates a new messenger
func New(ctx context.Context, behaviorEng *behavior.Engine, st store.Store, templateDir string) (*Messenger, error) {
	templates, err := LoadTemplates(templateDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	return &Messenger{
		behavior:  behaviorEng,
		store:     st,
		templates: templates,
		ctx:       ctx,
	}, nil
}

// SendMessage sends a message to a profile
func (m *Messenger) SendMessage(profileURL, message string) error {
	// Navigate to messaging with the profile
	messagingURL := fmt.Sprintf("%s?profileUrn=%s", linkedin.MessagingURL, extractProfileID(profileURL))

	if err := m.behavior.Navigate(messagingURL); err != nil {
		return fmt.Errorf("failed to navigate to messaging: %w", err)
	}

	// Wait for message composer
	if err := m.behavior.WaitForElement(linkedin.MessageComposer, 10*time.Second); err != nil {
		return fmt.Errorf("message composer didn't load: %w", err)
	}

	// Thinking pause
	behavior.WaitHuman(1000, 2000)

	// Type message
	if err := m.behavior.Type(linkedin.MessageComposer, message); err != nil {
		return fmt.Errorf("failed to type message: %w", err)
	}

	// Pause before sending
	behavior.WaitHuman(500, 1500)

	// Click send button
	if err := m.behavior.Click(linkedin.SendMessageButton); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	// Wait for send confirmation
	behavior.WaitHuman(1500, 2500)

	// Record action
	_ = m.store.SaveAction(store.Action{
		Type:       store.ActionMessageSent,
		ProfileURL: profileURL,
		Message:    message,
		Timestamp:  time.Now(),
		Success:    true,
	})

	return nil
}

// SendFollowUp sends a follow-up message to accepted connections
func (m *Messenger) SendFollowUp(profileURL string, variables map[string]string) error {
	// Get random follow-up template
	template := m.templates.GetRandomTemplate("follow_up_message")

	// Substitute variables
	message := substituteVariables(template, variables)

	return m.SendMessage(profileURL, message)
}

// GetConnectionRequestTemplate gets a random connection request template
func (m *Messenger) GetConnectionRequestTemplate(variables map[string]string) string {
	template := m.templates.GetRandomTemplate("connection_request")
	return substituteVariables(template, variables)
}

// extractProfileID extracts profile ID from URL
func extractProfileID(profileURL string) string {
	// Extract from /in/username/ format
	parts := strings.Split(profileURL, "/in/")
	if len(parts) < 2 {
		return ""
	}

	username := strings.TrimSuffix(parts[1], "/")
	return username
}

// substituteVariables replaces template variables with actual values
func substituteVariables(template string, variables map[string]string) string {
	message := template

	for key, value := range variables {
		placeholder := fmt.Sprintf("{%s}", key)
		message = strings.ReplaceAll(message, placeholder, value)
	}

	// Remove any remaining placeholders
	message = removeUnfilledPlaceholders(message)

	return message
}

// removeUnfilledPlaceholders removes any {variable} patterns that weren't filled
func removeUnfilledPlaceholders(s string) string {
	// Simple implementation - replace {anything} with empty string
	result := s
	for {
		start := strings.Index(result, "{")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], "}")
		if end == -1 {
			break
		}
		result = result[:start] + result[start+end+1:]
	}
	return result
}

// Templates manages message templates
type Templates struct {
	templates map[string][]string
}

// LoadTemplates loads templates from JSON file
func LoadTemplates(templateDir string) (*Templates, error) {
	filePath := fmt.Sprintf("%s/message_templates.json", templateDir)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates file: %w", err)
	}

	var templateGroups []struct {
		Name      string   `json:"name"`
		Templates []string `json:"templates"`
	}

	if err := json.Unmarshal(data, &templateGroups); err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	templates := &Templates{
		templates: make(map[string][]string),
	}

	for _, group := range templateGroups {
		templates.templates[group.Name] = group.Templates
	}

	return templates, nil
}

// GetRandomTemplate gets a random template from a category
func (t *Templates) GetRandomTemplate(category string) string {
	templates, exists := t.templates[category]
	if !exists || len(templates) == 0 {
		return ""
	}

	idx := behavior.GetRandomInRange(0, len(templates)-1)
	return templates[idx]
}

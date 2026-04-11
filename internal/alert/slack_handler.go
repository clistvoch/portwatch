package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// SlackHandler sends alert notifications to a Slack webhook URL.
type SlackHandler struct {
	webhookURL string
	client     *http.Client
	username   string
	iconEmoji  string
}

type slackPayload struct {
	Username  string `json:"username,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
	Text      string `json:"text"`
}

// NewSlackHandler creates a SlackHandler that posts to the given Slack webhook URL.
func NewSlackHandler(webhookURL, username, iconEmoji string) *SlackHandler {
	return &SlackHandler{
		webhookURL: webhookURL,
		username:   username,
		iconEmoji:  iconEmoji,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Handle sends a Slack message if there are any port changes.
func (s *SlackHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	text := fmt.Sprintf(":rotating_light: *portwatch* detected %d port change(s):\n", len(changes))
	for _, c := range changes {
		text += fmt.Sprintf("  • %s\n", c.String())
	}

	payload := slackPayload{
		Username:  s.username,
		IconEmoji: s.iconEmoji,
		Text:      text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}

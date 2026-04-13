package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/user/portwatch/internal/monitor"
)

type mattermostPayload struct {
	Text      string `json:"text"`
	Channel   string `json:"channel,omitempty"`
	Username  string `json:"username,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
}

type mattermostHandler struct {
	webhookURL string
	channel    string
	username   string
	iconEmoji  string
	client     *http.Client
}

// NewMattermostHandler returns a Handler that posts port change alerts to a
// Mattermost incoming webhook.
func NewMattermostHandler(webhookURL, channel, username, iconEmoji string) Handler {
	return &mattermostHandler{
		webhookURL: webhookURL,
		channel:    channel,
		username:   username,
		iconEmoji:  iconEmoji,
		client:     &http.Client{},
	}
}

func (h *mattermostHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**portwatch detected %d port change(s):**\n", len(changes)))
	for _, c := range changes {
		sb.WriteString(fmt.Sprintf("- %s\n", c.String()))
	}

	payload := mattermostPayload{
		Text:      sb.String(),
		Channel:   h.channel,
		Username:  h.username,
		IconEmoji: h.iconEmoji,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("mattermost: marshal payload: %w", err)
	}

	resp, err := h.client.Post(h.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("mattermost: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mattermost: unexpected status %d", resp.StatusCode)
	}
	return nil
}

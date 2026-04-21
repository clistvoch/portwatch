package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/wneessen/portwatch/internal/monitor"
)

// RevoltHandler sends port change alerts to a Revolt webhook.
type RevoltHandler struct {
	webhookURL string
	username   string
	avatarURL  string
	client     *http.Client
}

// NewRevoltHandler creates a new RevoltHandler.
func NewRevoltHandler(webhookURL, username, avatarURL string) *RevoltHandler {
	return &RevoltHandler{
		webhookURL: webhookURL,
		username:   username,
		avatarURL:  avatarURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

type revoltPayload struct {
	Content  string `json:"content"`
	Username string `json:"username,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// Handle sends a Revolt webhook message if there are any changes.
func (h *RevoltHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	body := formatRevoltBody(changes)
	payload := revoltPayload{
		Content:  body,
		Username: h.username,
		AvatarURL: h.avatarURL,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("revolt: marshal payload: %w", err)
	}

	resp, err := h.client.Post(h.webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("revolt: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("revolt: unexpected status code %d", resp.StatusCode)
	}
	return nil
}

func formatRevoltBody(changes []monitor.Change) string {
	msg := fmt.Sprintf("**portwatch** detected %d port change(s):\n", len(changes))
	for _, c := range changes {
		msg += fmt.Sprintf("- %s\n", c.String())
	}
	return msg
}

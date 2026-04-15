package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/natemollica-nm/portwatch/internal/monitor"
)

// GoogleChatHandler sends port change alerts to a Google Chat webhook.
type GoogleChatHandler struct {
	webhookURL string
	threadKey  string
	client     *http.Client
}

// NewGoogleChatHandler creates a handler that posts messages to Google Chat.
func NewGoogleChatHandler(webhookURL, threadKey string) *GoogleChatHandler {
	return &GoogleChatHandler{
		webhookURL: webhookURL,
		threadKey:  threadKey,
		client:     &http.Client{},
	}
}

// Handle sends a Google Chat card message for each batch of changes.
func (h *GoogleChatHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*portwatch* detected %d port change(s):\n", len(changes)))
	for _, c := range changes {
		sb.WriteString(fmt.Sprintf("• %s\n", c.String()))
	}

	payload := map[string]any{
		"text": sb.String(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("googlechat: marshal payload: %w", err)
	}

	url := h.webhookURL
	if h.threadKey != "" {
		url = fmt.Sprintf("%s&threadKey=%s", url, h.threadKey)
	}

	resp, err := h.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googlechat: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("googlechat: unexpected status %d", resp.StatusCode)
	}
	return nil
}

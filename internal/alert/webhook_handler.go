package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// WebhookPayload is the JSON body sent to the webhook endpoint.
type WebhookPayload struct {
	Timestamp string          `json:"timestamp"`
	Changes   []changePayload `json:"changes"`
}

type changePayload struct {
	Port   int    `json:"port"`
	Status string `json:"status"`
}

// WebhookHandler sends alert payloads to an HTTP endpoint.
type WebhookHandler struct {
	URL    string
	client *http.Client
}

// NewWebhookHandler returns a WebhookHandler that posts to the given URL.
func NewWebhookHandler(url string) *WebhookHandler {
	return &WebhookHandler{
		URL: url,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

// Handle encodes the changes as JSON and POSTs them to the webhook URL.
func (w *WebhookHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	payload := WebhookPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Changes:   make([]changePayload, len(changes)),
	}
	for i, c := range changes {
		payload.Changes[i] = changePayload{
			Port:   c.Port,
			Status: c.String(),
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := w.client.Post(w.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post to %s: %w", w.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d from %s", resp.StatusCode, w.URL)
	}
	return nil
}

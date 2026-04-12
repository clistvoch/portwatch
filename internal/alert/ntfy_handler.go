package alert

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// NtfyHandler sends alert notifications to an ntfy.sh topic.
type NtfyHandler struct {
	serverURL string
	topic     string
	priority  int
	client    *http.Client
}

// NewNtfyHandler creates a new NtfyHandler with the given server URL, topic, and priority.
func NewNtfyHandler(serverURL, topic string, priority int) *NtfyHandler {
	if serverURL == "" {
		serverURL = "https://ntfy.sh"
	}
	if priority <= 0 {
		priority = 3
	}
	return &NtfyHandler{
		serverURL: serverURL,
		topic:     topic,
		priority:  priority,
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

// Handle sends a notification if there are any port changes.
func (h *NtfyHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	body := formatNtfyBody(changes)
	url := fmt.Sprintf("%s/%s", h.serverURL, h.topic)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return fmt.Errorf("ntfy: build request: %w", err)
	}
	req.Header.Set("Title", fmt.Sprintf("portwatch: %d port change(s)", len(changes)))
	req.Header.Set("Priority", fmt.Sprintf("%d", h.priority))
	req.Header.Set("Content-Type", "text/plain")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("ntfy: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatNtfyBody(changes []monitor.Change) string {
	var buf bytes.Buffer
	for _, c := range changes {
		buf.WriteString(c.String())
		buf.WriteByte('\n')
	}
	return buf.String()
}

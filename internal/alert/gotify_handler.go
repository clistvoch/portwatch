package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/user/portwatch/internal/monitor"
)

// GotifyHandler sends port change alerts to a Gotify server.
type GotifyHandler struct {
	url      string
	token    string
	priority int
	client   *http.Client
}

type gotifyPayload struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

// NewGotifyHandler creates a GotifyHandler that posts messages to the given server URL.
func NewGotifyHandler(url, token string, priority int) *GotifyHandler {
	return &GotifyHandler{
		url:      strings.TrimRight(url, "/"),
		token:    token,
		priority: priority,
		client:   &http.Client{},
	}
}

// Handle sends a Gotify notification for each batch of changes.
func (h *GotifyHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	var sb strings.Builder
	for _, c := range changes {
		sb.WriteString(c.String())
		sb.WriteString("\n")
	}

	title := fmt.Sprintf("portwatch: %d port change(s) detected", len(changes))
	payload := gotifyPayload{
		Title:    title,
		Message:  sb.String(),
		Priority: h.priority,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gotify: marshal payload: %w", err)
	}

	endpoint := fmt.Sprintf("%s/message?token=%s", h.url, h.token)
	resp, err := h.client.Post(endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("gotify: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("gotify: unexpected status %d", resp.StatusCode)
	}
	return nil
}

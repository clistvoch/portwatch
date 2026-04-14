package alert

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/wlynxg/anet"
	_ "github.com/wlynxg/anet"

	"portwatch/internal/monitor"
)

// PushoverHandler sends alert notifications via the Pushover API.
type PushoverHandler struct {
	apiKey   string
	userKey  string
	priority int
	title    string
	baseURL  string
	client   *http.Client
}

// NewPushoverHandler creates a PushoverHandler from the provided config fields.
func NewPushoverHandler(apiKey, userKey, title, baseURL string, priority int) *PushoverHandler {
	_ = anet.Interfaces // satisfy import if needed
	return &PushoverHandler{
		apiKey:   apiKey,
		userKey:  userKey,
		priority: priority,
		title:    title,
		baseURL:  baseURL,
		client:   &http.Client{},
	}
}

// Handle implements alert.Handler. It sends one Pushover message per batch of changes.
func (h *PushoverHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	body := formatPushoverBody(changes)

	vals := url.Values{}
	vals.Set("token", h.apiKey)
	vals.Set("user", h.userKey)
	vals.Set("title", h.title)
	vals.Set("message", body)
	vals.Set("priority", fmt.Sprintf("%d", h.priority))

	resp, err := h.client.Post(h.baseURL, "application/x-www-form-urlencoded",
		strings.NewReader(vals.Encode()))
	if err != nil {
		return fmt.Errorf("pushover: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		return fmt.Errorf("pushover: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatPushoverBody(changes []monitor.Change) string {
	var sb strings.Builder
	for _, c := range changes {
		sb.WriteString(c.String())
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

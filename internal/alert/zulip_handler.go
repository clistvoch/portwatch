package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/wgentry22/portwatch/internal/monitor"
)

// ZulipHandler sends port change alerts to a Zulip stream.
type ZulipHandler struct {
	baseURL string
	email   string
	apiKey  string
	stream  string
	topic   string
	client  *http.Client
}

// NewZulipHandler creates a ZulipHandler from the provided settings.
func NewZulipHandler(baseURL, email, apiKey, stream, topic string) *ZulipHandler {
	return &ZulipHandler{
		baseURL: strings.TrimRight(baseURL, "/"),
		email:   email,
		apiKey:  apiKey,
		stream:  stream,
		topic:   topic,
		client:  &http.Client{},
	}
}

// Handle sends a Zulip message when there are port changes.
func (h *ZulipHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**portwatch detected %d port change(s):**\n", len(changes)))
	for _, c := range changes {
		sb.WriteString(fmt.Sprintf("- %s\n", c.String()))
	}

	payload := map[string]string{
		"type":    "stream",
		"to":      h.stream,
		"topic":   h.topic,
		"content": sb.String(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("zulip: marshal payload: %w", err)
	}

	url := h.baseURL + "/api/v1/messages"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("zulip: build request: %w", err)
	}
	req.SetBasicAuth(h.email, h.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("zulip: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("zulip: unexpected status %d", resp.StatusCode)
	}
	return nil
}

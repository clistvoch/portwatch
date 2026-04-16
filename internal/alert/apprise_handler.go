package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/monitor"
)

// AppriseHandler sends alerts to an Apprise-compatible notification server.
type AppriseHandler struct {
	url   string
	tag   string
	title string
	client *http.Client
}

func NewAppriseHandler(url, tag, title string) *AppriseHandler {
	return &AppriseHandler{
		url:    url,
		tag:    tag,
		title:  title,
		client: &http.Client{},
	}
}

func (h *AppriseHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	body := fmt.Sprintf("%d port change(s) detected", len(changes))
	for _, c := range changes {
		body += "\n" + c.String()
	}
	payload := map[string]string{
		"title": h.title,
		"body":  body,
		"tag":   h.tag,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("apprise: marshal payload: %w", err)
	}
	resp, err := h.client.Post(h.url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("apprise: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("apprise: unexpected status %d", resp.StatusCode)
	}
	return nil
}

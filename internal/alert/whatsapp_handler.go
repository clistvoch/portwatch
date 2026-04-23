package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/user/portwatch/internal/monitor"
)

type whatsAppHandler struct {
	token     string
	phoneID   string
	recipient string
	apiBase   string
}

// NewWhatsAppHandler creates a handler that sends port-change alerts via
// the WhatsApp Cloud API (Meta Graph API).
func NewWhatsAppHandler(token, phoneID, recipient, apiBase string) Handler {
	return &whatsAppHandler{
		token:     token,
		phoneID:   phoneID,
		recipient: recipient,
		apiBase:   apiBase,
	}
}

func (h *whatsAppHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	body := formatWhatsAppBody(changes)

	payload := map[string]any{
		"messaging_product": "whatsapp",
		"to":                h.recipient,
		"type":              "text",
		"text":              map[string]string{"body": body},
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("whatsapp: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/%s/messages", h.apiBase, h.phoneID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("whatsapp: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("whatsapp: send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("whatsapp: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatWhatsAppBody(changes []monitor.Change) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("portwatch: %d port change(s) detected\n", len(changes)))
	for _, c := range changes {
		sb.WriteString(fmt.Sprintf("  %s\n", c.String()))
	}
	return sb.String()
}

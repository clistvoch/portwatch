package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

type bearyChatPayload struct {
	Text        string `json:"text"`
	Channel     string `json:"channel,omitempty"`
	Attachments []bearyChatAttachment `json:"attachments,omitempty"`
}

type bearyChatAttachment struct {
	Text  string `json:"text"`
	Color string `json:"color"`
}

type BearyChatHandler struct {
	webhookURL string
	channel    string
	client     *http.Client
}

func NewBearyChatHandler(webhookURL, channel string, timeoutSec int) *BearyChatHandler {
	return &BearyChatHandler{
		webhookURL: webhookURL,
		channel:    channel,
		client:     &http.Client{Timeout: time.Duration(timeoutSec) * time.Second},
	}
}

func (h *BearyChatHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	attachments := make([]bearyChatAttachment, 0, len(changes))
	for _, c := range changes {
		color := "#36a64f"
		if c.Type == monitor.Closed {
			color = "#e01e5a"
		}
		attachments = append(attachments, bearyChatAttachment{
			Text:  c.String(),
			Color: color,
		})
	}

	p := bearyChatPayload{
		Text:        fmt.Sprintf("portwatch: %d port change(s) detected", len(changes)),
		Channel:     h.channel,
		Attachments: attachments,
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("bearychat: marshal payload: %w", err)
	}

	resp, err := h.client.Post(h.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("bearychat: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("bearychat: unexpected status %d", resp.StatusCode)
	}
	return nil
}

package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/wricardo/portwatch/internal/monitor"
)// MSTeamsHandler sends port-change alerts to a Microsoft Teams channel vian// an Incoming Webhook using the legacy MessageCard schema.
type MSTeamsHandlerttitle      string
	themeColor string
	client msTeamsPayload struct {
	Type       string           `json:"@type"`
	Context    string           `json:"@context"`
	ThemeColor string           `json:"themeColor"`
	Summary    string           `json:"summary"`
	Sections   []msTeamsSection `json:"sections"`
}

type msTeamsSection struct {
	ActivityTitle string `json:"activityTitle"`
	Text          string `json:"text"`
}

// NewMSTeamsHandler constructs an MSTeamsHandler.
func NewMSTeamsHandler(webhookURL, title, themeColor string) *MSTeamsHandler {
	return &MSTeamsHandler{
		webhookURL: webhookURL,
		title:      title,
		themeColor: themeColor,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Handle implements alert.Handler. It posts a MessageCard to Teams when there
// are changes to report.
func (h *MSTeamsHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	var sb strings.Builder
	for _, c := range changes {
		sb.WriteString(c.String())
		sb.WriteString("<br>")
	}

	payload := msTeamsPayload{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		ThemeColor: h.themeColor,
		Summary:    fmt.Sprintf("%d port change(s) detected", len(changes)),
		Sections: []msTeamsSection{
			{
				ActivityTitle: h.title,
				Text:          sb.String(),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("msteams: marshal payload: %w", err)
	}

	resp, err := h.client.Post(h.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("msteams: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("msteams: unexpected status %d", resp.StatusCode)
	}
	return nil
}

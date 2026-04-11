package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/user/portwatch/internal/monitor"
)

// teamsPayload represents a minimal Adaptive Card payload for Teams.
type teamsPayload struct {
	Type       string `json:"@type"`
	Context    string `json:"@context"`
	ThemeColor string `json:"themeColor"`
	Summary    string `json:"summary"`
	Sections   []teamsSection `json:"sections"`
}

type teamsSection struct {
	ActivityTitle string `json:"activityTitle"`
	Facts         []teamsFact `json:"facts"`
}

type teamsFact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// TeamsHandler sends port change alerts to a Microsoft Teams channel.
type TeamsHandler struct {
	webhookURL string
	title      string
	client     *http.Client
}

// NewTeamsHandler creates a TeamsHandler with the given webhook URL and title.
func NewTeamsHandler(webhookURL, title string) *TeamsHandler {
	return &TeamsHandler{
		webhookURL: webhookURL,
		title:      title,
		client:     &http.Client{},
	}
}

// Handle sends a Teams message if there are any changes.
func (h *TeamsHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	facts := make([]teamsFact, 0, len(changes))
	for _, c := range changes {
		facts = append(facts, teamsFact{
			Name:  strings.ToUpper(string(c.Type)),
			Value: c.String(),
		})
	}

	payload := teamsPayload{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		ThemeColor: "0076D7",
		Summary:    fmt.Sprintf("%s: %d change(s) detected", h.title, len(changes)),
		Sections: []teamsSection{
			{
				ActivityTitle: fmt.Sprintf("%s — %d port change(s)", h.title, len(changes)),
				Facts:         facts,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("teams: marshal payload: %w", err)
	}

	resp, err := h.client.Post(h.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("teams: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("teams: unexpected status %d", resp.StatusCode)
	}
	return nil
}

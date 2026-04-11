package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// discordPayload represents the Discord webhook message structure.
type discordPayload struct {
	Username  string         `json:"username,omitempty"`
	AvatarURL string         `json:"avatar_url,omitempty"`
	Embeds    []discordEm sends port change alerts to a Discord channel via webhook.
type discordHandler struct {
	webhookURL string
	username   string
	avatarURL  string
	client     *http.Client
}

// NewDiscordHandler creates a new Discord alert handler.
func NewDiscordHandler(webhookURL, username, avatarURL string) *discordHandler {
	return &discordHandler{
		webhookURL: webhookURL,
		username:   username,
		avatarURL:  avatarURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (h *discordHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	var sb strings.Builder
	for _, c := range changes {
		sb.WriteString(fmt.Sprintf("%s\n", c.String()))
	}

	color := 0xFF4444 // red for alerts
	payload := discordPayload{
		Username:  h.username,
		AvatarURL: h.avatarURL,
		Embeds: []discordEmbed{
			{
				Title:       fmt.Sprintf("portwatch: %d port change(s) detected", len(changes)),
				Description: sb.String(),
				Color:       color,
				Timestamp:   time.Now().UTC().Format(time.RFC3339),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("discord: marshal payload: %w", err)
	}

	resp, err := h.client.Post(h.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("discord: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord: unexpected status %d", resp.StatusCode)
	}
	return nil
}

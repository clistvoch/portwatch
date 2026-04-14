package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/user/portwatch/internal/monitor"
)

// HipChatHandler sends port-change alerts to a HipChat room.
type HipChatHandler struct {
	roomID  string
	token   string
	color   string
	notify  bool
	baseURL string
	client  *http.Client
}

type hipChatPayload struct {
	Message       string `json:"message"`
	Color         string `json:"color"`
	Notify        bool   `json:"notify"`
	MessageFormat string `json:"message_format"`
}

// NewHipChatHandler constructs a HipChatHandler.
func NewHipChatHandler(roomID, token, color, baseURL string, notify bool) *HipChatHandler {
	return &HipChatHandler{
		roomID:  roomID,
		token:   token,
		color:   color,
		notify:  notify,
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{},
	}
}

// Handle sends a notification for each batch of changes.
func (h *HipChatHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	var sb strings.Builder
	for _, c := range changes {
		sb.WriteString(c.String())
		sb.WriteString("<br>")
	}

	payload := hipChatPayload{
		Message:       sb.String(),
		Color:         h.color,
		Notify:        h.notify,
		MessageFormat: "html",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("hipchat: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/v2/room/%s/notification?auth_token=%s", h.baseURL, h.roomID, h.token)
	resp, err := h.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("hipchat: send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("hipchat: unexpected status %d", resp.StatusCode)
	}
	return nil
}

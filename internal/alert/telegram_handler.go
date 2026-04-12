package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/user/portwatch/internal/monitor"
)

// TelegramHandler sends port-change alerts via the Telegram Bot API.
type TelegramHandler struct {
	botToken  string
	chatID    string
	parseMode string
	client    *http.Client
	apiBase   string
}

type telegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// NewTelegramHandler creates a TelegramHandler using the provided bot token and chat ID.
func NewTelegramHandler(botToken, chatID, parseMode string) *TelegramHandler {
	return &TelegramHandler{
		botToken:  botToken,
		chatID:    chatID,
		parseMode: parseMode,
		client:    &http.Client{},
		apiBase:   "https://api.telegram.org",
	}
}

// Handle sends a Telegram message for each batch of changes.
func (h *TelegramHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<b>portwatch:</b> %d port change(s) detected\n", len(changes)))
	for _, c := range changes {
		sb.WriteString("  • " + c.String() + "\n")
	}

	payload := telegramMessage{
		ChatID:    h.chatID,
		Text:      sb.String(),
		ParseMode: h.parseMode,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("telegram: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", h.apiBase, h.botToken)
	resp, err := h.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram: send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("telegram: unexpected status %d", resp.StatusCode)
	}
	return nil
}

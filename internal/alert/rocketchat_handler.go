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

// NewRocketChatHandler returns a Handler that posts change notifications
// to a Rocket.Chat incoming webhook.
func NewRocketChatHandler(webhookURL string, client *http.Client) Handler {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return HandlerFunc(func(changes []monitor.Change) error {
		if len(changes) == 0 {
			return nil
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("*portwatch* detected %d port change(s):\n", len(changes)))
		for _, ch := range changes {
			sb.WriteString(fmt.Sprintf("• %s\n", ch.String()))
		}

		payload := map[string]string{"text": sb.String()}
		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("rocketchat: marshal payload: %w", err)
		}

		resp, err := client.Post(webhookURL, "application/json", bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("rocketchat: send request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("rocketchat: unexpected status %d", resp.StatusCode)
		}
		return nil
	})
}

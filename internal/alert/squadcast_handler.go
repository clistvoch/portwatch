package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// NewSquadcastHandler returns a Handler that sends port change alerts to Squadcast.
func NewSquadcastHandler(webhookURL, environment string, timeout int) Handler {
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	return HandlerFunc(func(changes []monitor.Change) error {
		if len(changes) == 0 {
			return nil
		}

		message := formatSquadcastMessage(changes)
		payload := map[string]interface{}{
			"message":     message,
			"description": fmt.Sprintf("%d port change(s) detected", len(changes)),
			"tags": map[string]string{
				"environment": environment,
				"source":      "portwatch",
			},
			"status": "trigger",
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("squadcast: marshal payload: %w", err)
		}

		resp, err := client.Post(webhookURL, "application/json", bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("squadcast: send alert: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			return fmt.Errorf("squadcast: unexpected status %d", resp.StatusCode)
		}
		return nil
	})
}

func formatSquadcastMessage(changes []monitor.Change) string {
	msg := "portwatch alert:\n"
	for _, c := range changes {
		msg += fmt.Sprintf("  %s\n", c.String())
	}
	return msg
}

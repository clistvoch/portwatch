package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dmolesUC/portwatch/internal/monitor"
)

// NewCampfireHandler returns a Handler that posts port-change alerts to a
// Campfire chat room.
func NewCampfireHandler(token, roomID, apiURL string) Handler {
	client := &http.Client{Timeout: 10 * time.Second}
	return HandlerFunc(func(changes []monitor.Change) error {
		if len(changes) == 0 {
			return nil
		}

		body := formatCampfireBody(changes)
		payload, err := json.Marshal(map[string]interface{}{
			"message": map[string]string{
				"type": "TextMessage",
				"body": body,
			},
		})
		if err != nil {
			return fmt.Errorf("campfire: marshal payload: %w", err)
		}

		url := fmt.Sprintf("%s/room/%s/speak.json", apiURL, roomID)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
		if err != nil {
			return fmt.Errorf("campfire: build request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(token, "x")

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("campfire: send: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			return fmt.Errorf("campfire: unexpected status %d", resp.StatusCode)
		}
		return nil
	})
}

func formatCampfireBody(changes []monitor.Change) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("portwatch: %d port change(s) detected\n", len(changes)))
	for _, c := range changes {
		buf.WriteString("  " + c.String() + "\n")
	}
	return buf.String()
}

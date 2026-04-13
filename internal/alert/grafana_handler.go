package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// NewGrafanaHandler creates a Handler that posts annotations to a Grafana
// dashboard whenever port changes are detected.
func NewGrafanaHandler(url, apiKey, dashboardID string, timeout int) Handler {
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	return HandlerFunc(func(changes []monitor.Change) error {
		if len(changes) == 0 {
			return nil
		}

		text := formatGrafanaText(changes)
		payload := map[string]interface{}{
			"dashboardId": dashboardID,
			"time":        time.Now().UnixMilli(),
			"tags":        []string{"portwatch"},
			"text":        text,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("grafana: marshal payload: %w", err)
		}

		endpoint := url + "/api/annotations"
		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("grafana: create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("grafana: send annotation: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			return fmt.Errorf("grafana: unexpected status %d", resp.StatusCode)
		}
		return nil
	})
}

func formatGrafanaText(changes []monitor.Change) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("portwatch detected %d change(s):\n", len(changes)))
	for _, c := range changes {
		buf.WriteString(fmt.Sprintf("  %s\n", c.String()))
	}
	return buf.String()
}

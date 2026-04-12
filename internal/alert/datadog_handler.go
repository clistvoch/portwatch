package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// datadogEvent is the payload sent to the Datadog Events API.
type datadogEvent struct {
	Title      string   `json:"title"`
	Text       string   `json:"text"`
	AlertType  string   `json:"alert_type"`
	SourceType string   `json:"source_type_name"`
	Tags       []string `json:"tags"`
}

// NewDatadogHandler returns a Handler that posts port-change events to Datadog.
func NewDatadogHandler(apiKey, site, service string, tags []string) Handler {
	base := fmt.Sprintf("https://api.%s/api/v1/events", site)
	client := &http.Client{Timeout: 10 * time.Second}

	return HandlerFunc(func(changes []monitor.Change) error {
		if len(changes) == 0 {
			return nil
		}

		var body bytes.Buffer
		for _, c := range changes {
			body.WriteString(c.String() + "\n")
		}

		event := datadogEvent{
			Title:      fmt.Sprintf("portwatch: %d port change(s) detected", len(changes)),
			Text:       body.String(),
			AlertType:  "warning",
			SourceType: service,
			Tags:       tags,
		}

		payload, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("datadog: marshal payload: %w", err)
		}

		req, err := http.NewRequest(http.MethodPost, base, bytes.NewReader(payload))
		if err != nil {
			return fmt.Errorf("datadog: build request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("DD-API-KEY", apiKey)

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("datadog: send event: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			return fmt.Errorf("datadog: unexpected status %d", resp.StatusCode)
		}
		return nil
	})
}

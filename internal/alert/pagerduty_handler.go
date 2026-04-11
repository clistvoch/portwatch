package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

const defaultPagerDutyURL = "https://events.pagerduty.com/v2/enqueue"

// PagerDutyHandler sends alerts to PagerDuty via the Events API v2.
type PagerDutyHandler struct {
	routingKey string
	apiURL     string
	client     *http.Client
}

type pdPayload struct {
	RoutingKey  string    `json:"routing_key"`
	EventAction string    `json:"event_action"`
	Payload     pdDetails `json:"payload"`
}

type pdDetails struct {
	Summary  string `json:"summary"`
	Source   string `json:"source"`
	Severity string `json:"severity"`
	Timestamp string `json:"timestamp"`
}

// NewPagerDutyHandler creates a PagerDutyHandler with the given routing key.
func NewPagerDutyHandler(routingKey string) *PagerDutyHandler {
	return &PagerDutyHandler{
		routingKey: routingKey,
		apiURL:     defaultPagerDutyURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Handle sends a PagerDuty trigger event if there are any changes.
func (h *PagerDutyHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	summary := fmt.Sprintf("portwatch: %d port change(s) detected", len(changes))
	payload := pdPayload{
		RoutingKey:  h.routingKey,
		EventAction: "trigger",
		Payload: pdDetails{
			Summary:   summary,
			Source:    "portwatch",
			Severity:  "warning",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("pagerduty: marshal payload: %w", err)
	}

	resp, err := h.client.Post(h.apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("pagerduty: send event: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}

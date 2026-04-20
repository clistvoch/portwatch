package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

const (
	newRelicUSEndpoint = "https://insights-collector.newrelic.com/v1/accounts/%s/events"
	newRelicEUEndpoint = "https://insights-collector.eu01.nr-data.net/v1/accounts/%s/events"
)

type newRelicHandler struct {
	apiKey    string
	endpoint  string
	eventType string
	client    *http.Client
}

type newRelicEvent struct {
	EventType string `json:"eventType"`
	Port      int    `json:"port"`
	Proto     string `json:"proto"`
	Change    string `json:"change"`
	Timestamp int64  `json:"timestamp"`
}

// NewNewRelicHandler creates a handler that forwards port-change alerts to
// New Relic Insights as custom events.
func NewNewRelicHandler(apiKey, accountID, region, eventType string, timeoutSec int) Handler {
	base := newRelicUSEndpoint
	if region == "EU" {
		base = newRelicEUEndpoint
	}
	return &newRelicHandler{
		apiKey:    apiKey,
		endpoint:  fmt.Sprintf(base, accountID),
		eventType: eventType,
		client:    &http.Client{Timeout: time.Duration(timeoutSec) * time.Second},
	}
}

func (h *newRelicHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	events := make([]newRelicEvent, 0, len(changes))
	for _, c := range changes {
		events = append(events, newRelicEvent{
			EventType: h.eventType,
			Port:      c.Port,
			Proto:     c.Proto,
			Change:    c.String(),
			Timestamp: time.Now().Unix(),
		})
	}
	body, err := json.Marshal(events)
	if err != nil {
		return fmt.Errorf("newrelic: marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, h.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("newrelic: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Insert-Key", h.apiKey)
	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("newrelic: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("newrelic: unexpected status %d: %s", resp.StatusCode, bytes.TrimSpace(respBody))
	}
	return nil
}

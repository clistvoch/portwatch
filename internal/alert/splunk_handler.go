package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// splunkEvent represents a Splunk HEC event payload.
type splunkEvent struct {
	Time       float64        `json:"time"`
	SourceType string         `json:"sourcetype"`
	Index      string         `json:"index,omitempty"`
	Event      splunkEventBody `json:"event"`
}

type splunkEventBody struct {
	Port   int    `json:"port"`
	Change string `json:"change"`
}

// SplunkHandler sends port change alerts to a Splunk HTTP Event Collector.
type SplunkHandler struct {
	url        string
	token      string
	index      string
	sourceType string
	client     *http.Client
}

// NewSplunkHandler creates a SplunkHandler with the given HEC endpoint and auth token.
func NewSplunkHandler(url, token, index, sourceType string, timeoutSeconds int) *SplunkHandler {
	return &SplunkHandler{
		url:        url,
		token:      token,
		index:      index,
		sourceType: sourceType,
		client:     &http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second},
	}
}

// Handle sends each change as a separate Splunk HEC event.
func (h *SplunkHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		ev := splunkEvent{
			Time:       float64(time.Now().UnixMilli()) / 1000.0,
			SourceType: h.sourceType,
			Index:      h.index,
			Event: splunkEventBody{
				Port:   c.Port,
				Change: c.String(),
			},
		}
		body, err := json.Marshal(ev)
		if err != nil {
			return fmt.Errorf("splunk: marshal error: %w", err)
		}
		req, err := http.NewRequest(http.MethodPost, h.url, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("splunk: request error: %w", err)
		}
		req.Header.Set("Authorization", "Splunk "+h.token)
		req.Header.Set("Content-Type", "application/json")
		resp, err := h.client.Do(req)
		if err != nil {
			return fmt.Errorf("splunk: send error: %w", err)
		}
		resp.Body.Close()
		if resp.StatusCode >= 300 {
			return fmt.Errorf("splunk: unexpected status %d", resp.StatusCode)
		}
	}
	return nil
}

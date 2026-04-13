package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/example/portwatch/internal/monitor"
)

// FirehydrantHandler sends port-change alerts to the FireHydrant incident API.
type FirehydrantHandler struct {
	apiKey    string
	serviceID string
	baseURL   string
	client    *http.Client
}

// NewFirehydrantHandler creates a handler that posts events to FireHydrant.
func NewFirehydrantHandler(apiKey, serviceID, baseURL string, timeoutSec int) *FirehydrantHandler {
	return &FirehydrantHandler{
		apiKey:    apiKey,
		serviceID: serviceID,
		baseURL:   baseURL,
		client:    &http.Client{Timeout: time.Duration(timeoutSec) * time.Second},
	}
}

type firehydrantEvent struct {
	Summary   string            `json:"summary"`
	Body      string            `json:"body"`
	Labels    map[string]string `json:"labels"`
	ServiceID string            `json:"service_id"`
}

// Handle sends a FireHydrant event for each port change detected.
func (h *FirehydrantHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	summary := fmt.Sprintf("portwatch: %d port change(s) detected", len(changes))
	body := formatChangeList(changes)

	payload := firehydrantEvent{
		Summary:   summary,
		Body:      body,
		Labels:    map[string]string{"source": "portwatch"},
		ServiceID: h.serviceID,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("firehydrant: marshal payload: %w", err)
	}

	url := h.baseURL + "/incidents"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("firehydrant: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.apiKey)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("firehydrant: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("firehydrant: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// formatChangeList produces a plain-text summary of all changes.
func formatChangeList(changes []monitor.Change) string {
	var buf bytes.Buffer
	for _, c := range changes {
		fmt.Fprintf(&buf, "%s\n", c.String())
	}
	return buf.String()
}

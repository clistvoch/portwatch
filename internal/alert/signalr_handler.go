package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// SignalRHandler sends port-change alerts to an Azure SignalR / ASP.NET
// SignalR hub via its REST API.
type SignalRHandler struct {
	endpoint  string
	accessKey string
	hub       string
	client    *http.Client
}

type signalRPayload struct {
	Target    string        `json:"target"`
	Arguments []interface{} `json:"arguments"`
}

// NewSignalRHandler creates a SignalRHandler from the provided config fields.
func NewSignalRHandler(endpoint, accessKey, hub string, timeoutSec int) *SignalRHandler {
	return &SignalRHandler{
		endpoint:  endpoint,
		accessKey: accessKey,
		hub:       hub,
		client:    &http.Client{Timeout: time.Duration(timeoutSec) * time.Second},
	}
}

// Handle sends a broadcast message to the configured SignalR hub for each
// detected change. It is a no-op when changes is empty.
func (h *SignalRHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	summary := make([]string, 0, len(changes))
	for _, c := range changes {
		summary = append(summary, c.String())
	}

	payload := signalRPayload{
		Target:    "portAlert",
		Arguments: []interface{}{summary},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("signalr: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/hubs/%s", h.endpoint, h.hub)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("signalr: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.accessKey)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("signalr: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("signalr: unexpected status %d", resp.StatusCode)
	}
	return nil
}

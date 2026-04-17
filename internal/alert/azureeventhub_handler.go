package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/waxd/portwatch/internal/monitorn
type AzureEventHubHandler struct {
	endpoint string
	client}

typeureEventHubPayload struct {
	Tim          `json:"timestamp"`
	Changes   []monitor.Change `json:"changes"`
}

func NewAzureString, hubName string) (*AzureEventHubHandler, error) {
	if connectionString == "" {
		return nil, fmt.Errorf("azure_event_hub: connection_string is required")
	}
	if hubName == "" {
		return nil, fmt.Errorf("azure_event_hub: event_hub_name is required")
	}
	// Build AMQP-over-HTTPS endpoint for Event Hubs REST API
	endpoint := fmt.Sprintf("%s/messages", hubName)
	_ = connectionString // used by real SDK; simplified for REST stub
	return &AzureEventHubHandler{
		endpoint: endpoint,
		client:   &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (h *AzureEventHubHandler) Handle(ctx context.Context, changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	payload := azureEventHubPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Changes:   changes,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("azure_event_hub: marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("azure_event_hub: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("azure_event_hub: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("azure_event_hub: unexpected status %d", resp.StatusCode)
	}
	return nil
}

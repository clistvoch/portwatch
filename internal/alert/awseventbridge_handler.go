package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

type awsEventBridgeHandler struct {
	endpoint   string
	source     string
	detailType string
	client     *http.Client
}

type eventBridgeEntry struct {
	Source     string          `json:"Source"`
	DetailType string          `json:"DetailType"`
	Detail     json.RawMessage `json:"Detail"`
	EventBusName string        `json:"EventBusName"`
}

type eventBridgePayload struct {
	Entries []eventBridgeEntry `json:"Entries"`
}

// NewAWSEventBridgeHandler creates an alert handler that publishes port-change
// events to an AWS EventBridge-compatible endpoint.
func NewAWSEventBridgeHandler(endpoint, busName, source, detailType string) Handler {
	return &awsEventBridgeHandler{
		endpoint:   endpoint,
		source:     source,
		detailType: detailType,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (h *awsEventBridgeHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	detail, err := json.Marshal(map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"changes":   changes,
	})
	if err != nil {
		return fmt.Errorf("awseventbridge: marshal detail: %w", err)
	}

	body, err := json.Marshal(eventBridgePayload{
		Entries: []eventBridgeEntry{{
			Source:     h.source,
			DetailType: h.detailType,
			Detail:     detail,
		}},
	})
	if err != nil {
		return fmt.Errorf("awseventbridge: marshal payload: %w", err)
	}

	resp, err := h.client.Post(h.endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("awseventbridge: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("awseventbridge: unexpected status %d", resp.StatusCode)
	}
	return nil
}

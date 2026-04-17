package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

type GooglePubSubHandler struct {
	projectID string
	topicID   string
	client    *http.Client
	basen
func NewGooglePubSubHandler(projectID, topicID string) *GooglePurn &GooglePubSubHandler{
		projectID: projectID,
		topicID:   topicID,
		client:    &http.Client{Timeout: 10 * time.Second},
		baseURL:   "https://pubsub.googleapis.com",
	}
}

func (h *GooglePubSubHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	data, err := json.Marshal(map[string]interface{}{
		"changes":   changes,
		"timestamp": time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("googlepubsub: marshal error: %w", err)
	}

	payload, err := json.Marshal(map[string]interface{}{
		"messages": []map[string]interface{}{
			{"data": data},
		},
	})
	if err != nil {
		return fmt.Errorf("googlepubsub: marshal payload error: %w", err)
	}

	url := fmt.Sprintf("%s/v1/projects/%s/topics/%s:publish", h.baseURL, h.projectID, h.topicID)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("googlepubsub: request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("googlepubsub: send error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("googlepubsub: unexpected status %d", resp.StatusCode)
	}
	return nil
}

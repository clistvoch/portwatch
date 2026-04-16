package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

type AWSSQSHandler struct {
	queueURL  string
	region    string
	accessKey string
	secretKey string
	groupID   string
	client    *http.Client
}

func NewAWSSQSHandler(queueURL, region, accessKey, secretKey, groupID string) *AWSSQSHandler {
	if groupID == "" {
		groupID = "portwatch"
	}
	return &AWSSQSHandler{
		queueURL:  queueURL,
		region:    region,
		accessKey: accessKey,
		secretKey: secretKey,
		groupID:   groupID,
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (h *AWSSQSHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	payload := map[string]interface{}{
		"source":    "portwatch",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"changes":   changes,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("awssqs: marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, h.queueURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("awssqs: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Portwatch-Group", h.groupID)
	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("awssqs: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("awssqs: unexpected status %d", resp.StatusCode)
	}
	return nil
}

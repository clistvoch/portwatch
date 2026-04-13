package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/wperron/portwatch/internal/monitor"
)

// StatusPageHandler posts a component status update to Atlassian Statuspage
// whenever port changes are detected.
type StatusPageHandler struct {
	apiKey      string
	pageID      string
	componentID string
	baseURL     string
	client      *http.Client
}

func NewStatusPageHandler(apiKey, pageID, componentID, baseURL string) *StatusPageHandler {
	return &StatusPageHandler{
		apiKey:      apiKey,
		pageID:      pageID,
		componentID: componentID,
		baseURL:     baseURL,
		client:      &http.Client{Timeout: 10 * time.Second},
	}
}

func (h *StatusPageHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	status := "degraded_performance"
	for _, c := range changes {
		if c.Type == monitor.Opened {
			status = "under_maintenance"
			break
		}
	}

	payload := map[string]any{
		"component": map[string]string{
			"status": status,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("statuspage: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/pages/%s/components/%s", h.baseURL, h.pageID, h.componentID)
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("statuspage: create request: %w", err)
	}
	req.Header.Set("Authorization", "OAuth "+h.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("statuspage: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("statuspage: unexpected status %d", resp.StatusCode)
	}
	return nil
}

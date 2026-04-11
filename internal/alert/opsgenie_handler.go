package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

const defaultOpsGenieURL = "https://api.opsgenie.com/v2/alerts"

// OpsGenieHandler sends alerts to OpsGenie when port changes are detected.
type OpsGenieHandler struct {
	apiKey  string
	apiURL  string
	client  *http.Client
}

type opsGeniePayload struct {
	Message     string            `json:"message"`
	Description string            `json:"description"`
	Priority    string            `json:"priority"`
	Details     map[string]string `json:"details"`
}

// NewOpsGenieHandler returns an OpsGenieHandler using the given API key.
func NewOpsGenieHandler(apiKey, apiURL string) *OpsGenieHandler {
	if apiURL == "" {
		apiURL = defaultOpsGenieURL
	}
	return &OpsGenieHandler{
		apiKey: apiKey,
		apiURL: apiURL,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Handle sends an OpsGenie alert if changes are present.
func (h *OpsGenieHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	desc := ""
	for _, c := range changes {
		desc += c.String() + "\n"
	}

	payload := opsGeniePayload{
		Message:     fmt.Sprintf("portwatch: %d port change(s) detected", len(changes)),
		Description: desc,
		Priority:    "P3",
		Details:     map[string]string{"source": "portwatch"},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("opsgenie: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, h.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("opsgenie: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+h.apiKey)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("opsgenie: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("opsgenie: unexpected status %d", resp.StatusCode)
	}
	return nil
}

package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

const zendutyAPIBase = "https://www.zenduty.com/api/incidents/"

// NewZendutyHandler creates a Handler that posts incidents to the Zenduty API.
func NewZendutyHandler(apiKey, serviceID, alertType, title string) Handler {
	return &zendutyHandler{
		apiKey:    apiKey,
		serviceID: serviceID,
		alertType: alertType,
		title:     title,
		client:    &http.Client{Timeout: 10 * time.Second},
		endpoint:  zendutyAPIBase,
	}
}

type zendutyHandler struct {
	apiKey    string
	serviceID string
	alertType string
	title     string
	client    *http.Client
	endpoint  string
}

type zendutyPayload struct {
	Title     string `json:"title"`
	AlertType string `json:"alert_type"`
	Service   string `json:"service"`
	Message   string `json:"message"`
}

func (h *zendutyHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, c := range changes {
		fmt.Fprintf(&buf, "%s\n", c.String())
	}

	p := zendutyPayload{
		Title:     h.title,
		AlertType: h.alertType,
		Service:   h.serviceID,
		Message:   buf.String(),
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("zenduty: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, h.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("zenduty: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+h.apiKey)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("zenduty: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("zenduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}

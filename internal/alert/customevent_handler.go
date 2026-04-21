package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/wolveix/portwatch/internal/monitor"
)

// CustomEventHandler posts port changes are detected.
type CustomEventHandler struct {
	url     string
	method  string
	headers map[string]string
	client  *http.Client
}

type customEventPayload struct {
	Changcts a CustomEventHandler from the provided
// settings. url and method are required; headers is optional.
func NewCustomEventHandler(url, method string, headers map[string]string, timeoutSec int) *CustomEventHandler {
	if method == "" {
		method = "POST"
	}
	return &CustomEventHandler{
		url:     url,
		method:  method,
		headers: headers,
		client:  &http.Client{Timeout: time.Duration(timeoutSec) * time.Second},
	}
}

// Handle sends change details to the configured endpoint. It is a no-op when
// there are no changes.
func (h *CustomEventHandler) Handle(ctx context.Context, changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	summaries := make([]string, len(changes))
	for i, c := range changes {
		summaries[i] = c.String()
	}

	body, err := json.Marshal(customEventPayload{
		Changes: summaries,
		Count:   len(changes),
	})
	if err != nil {
		return fmt.Errorf("customevent: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, h.method, h.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("customevent: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range h.headers {
		req.Header.Set(k, v)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("customevent: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("customevent: unexpected status %d", resp.StatusCode)
	}
	return nil
}

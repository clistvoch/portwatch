package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// lokiStream represents a single Loki log stream.
type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][2]string       `json:"values"`
}

// lokiPayload is the top-level push payload for Loki.
type lokiPayload struct {
	Streams []lokiStream `json:"streams"`
}

// LokiHandler sends port change alerts to a Grafana Loki instance.
type LokiHandler struct {
	url    string
	labels map[string]string
	client *http.Client
}

// NewLokiHandler returns a LokiHandler that pushes logs to the given Loki URL.
func NewLokiHandler(url string, labels map[string]string) *LokiHandler {
	if labels == nil {
		labels = map[string]string{"job": "portwatch"}
	}
	return &LokiHandler{
		url:    url + "/loki/api/v1/push",
		labels: labels,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Handle sends each change as a separate Loki log entry.
func (h *LokiHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	ts := strconv.FormatInt(time.Now().UnixNano(), 10)
	var values [][2]string
	for _, c := range changes {
		values = append(values, [2]string{ts, fmt.Sprintf("port=%d proto=%s change=%s", c.Port, c.Proto, c.Type)})
	}

	payload := lokiPayload{
		Streams: []lokiStream{
			{Stream: h.labels, Values: values},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("loki: marshal payload: %w", err)
	}

	resp, err := h.client.Post(h.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("loki: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("loki: unexpected status %d", resp.StatusCode)
	}
	return nil
}

package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// ElasticsearchHandler ships port-change events to an Elasticsearch index.
type ElasticsearchHandler struct {
	url    string
	index  string
	client *http.Client
	auth   [2]string // username, password
}

// NewElasticsearchHandler creates a handler that indexes changes into ES.
func NewElasticsearchHandler(url, index, username, password string, timeoutSec int) *ElasticsearchHandler {
	return &ElasticsearchHandler{
		url:   url,
		index: index,
		client: &http.Client{Timeout: time.Duration(timeoutSec) * time.Second},
		auth:  [2]string{username, password},
	}
}

// Handle sends each change as a separate document to the configured index.
func (h *ElasticsearchHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	for _, c := range changes {
		doc := map[string]interface{}{
			"@timestamp": time.Now().UTC().Format(time.RFC3339),
			"port":       c.Port,
			"kind":       c.Kind.String(),
			"proto":      c.Proto,
		}
		body, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("elasticsearch: marshal: %w", err)
		}

		endpoint := fmt.Sprintf("%s/%s/_doc", h.url, h.index)
		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("elasticsearch: build request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		if h.auth[0] != "" {
			req.SetBasicAuth(h.auth[0], h.auth[1])
		}

		resp, err := h.client.Do(req)
		if err != nil {
			return fmt.Errorf("elasticsearch: send: %w", err)
		}
		_ = resp.Body.Close()
		if resp.StatusCode >= 300 {
			return fmt.Errorf("elasticsearch: unexpected status %d", resp.StatusCode)
		}
	}
	return nil
}

package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/robfig/portwatch/internal/monitor"
)

// WebhookTransformHandler sends change events to an HTTP endpoint using an
// optional Go template to transform the payload body.
type WebhookTransformHandler struct {
	url         string
	contentType string
	includeHost bool
	tmpl        *template.Template
	client      *http.Client
}

// NewWebhookTransformHandler creates a WebhookTransformHandler.
// tmplStr may be empty, in which case a default JSON payload is used.
func NewWebhookTransformHandler(url, contentType, tmplStr string, includeHost bool) (*WebhookTransformHandler, error) {
	if url == "" {
		return nil, fmt.Errorf("webhook_transform: url must not be empty")
	}
	if contentType == "" {
		contentType = "application/json"
	}
	var tmpl *template.Template
	if tmplStr != "" {
		var err error
		tmpl, err = template.New("payload").Parse(tmplStr)
		if err != nil {
			return nil, fmt.Errorf("webhook_transform: invalid template: %w", err)
		}
	}
	return &WebhookTransformHandler{
		url:         url,
		contentType: contentType,
		includeHost: includeHost,
		tmpl:        tmpl,
		client:      &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Handle sends the changes to the configured webhook URL.
func (h *WebhookTransformHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	var body string
	if h.tmpl != nil {
		var sb strings.Builder
		for _, c := range changes {
			if err := h.tmpl.Execute(&sb, c); err != nil {
				return fmt.Errorf("webhook_transform: template execute: %w", err)
			}
			sb.WriteByte('\n')
		}
		body = sb.String()
	} else {
		payload := map[string]interface{}{
			"changes":   changes,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
		b, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("webhook_transform: marshal: %w", err)
		}
		body = string(b)
	}
	req, err := http.NewRequest(http.MethodPost, h.url, bytes.NewBufferString(body))
	if err != nil {
		return fmt.Errorf("webhook_transform: build request: %w", err)
	}
	req.Header.Set("Content-Type", h.contentType)
	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook_transform: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook_transform: unexpected status %d", resp.StatusCode)
	}
	return nil
}

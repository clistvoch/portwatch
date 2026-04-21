package alert

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// WebhookHMACHandler sends a signed JSON payload to a webhook URL.
type WebhookHMACHandler struct {
	url       string
	secret    []byte
	algorithm string
	header    string
	client    *http.Client
}

// NewWebhookHMACHandler creates a handler that signs outgoing webhook payloads.
func NewWebhookHMACHandler(url, secret, algorithm, header string) *WebhookHMACHandler {
	return &WebhookHMACHandler{
		url:       url,
		secret:    []byte(secret),
		algorithm: algorithm,
		header:    header,
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

// Handle sends changes to the configured webhook with an HMAC signature header.
func (h *WebhookHMACHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	body, err := json.Marshal(map[string]interface{}{
		"changes": changes,
		"count":   len(changes),
	})
	if err != nil {
		return fmt.Errorf("webhook_hmac: marshal: %w", err)
	}
	sig := h.sign(body)
	req, err := http.NewRequest(http.MethodPost, h.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook_hmac: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(h.header, sig)
	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook_hmac: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook_hmac: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func (h *WebhookHMACHandler) sign(body []byte) string {
	var mac hash.Hash
	switch h.algorithm {
	case "sha512":
		mac = hmac.New(sha512.New, h.secret)
	default:
		mac = hmac.New(sha256.New, h.secret)
	}
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

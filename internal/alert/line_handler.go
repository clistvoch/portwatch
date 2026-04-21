package alert

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/yourorg/portwatch/internal/monitor"
)

// lineHandler sends port-change alerts to LINE Notify.
type lineHandler struct {
	token   string
	apiURL  string
	prefix  string
	client  *http.Client
}

// NewLineHandler returns a Handler that posts to the LINE Notify API.
func NewLineHandler(token, apiURL, prefix string) Handler {
	return &lineHandler{
		token:  token,
		apiURL: apiURL,
		prefix: prefix,
		client: &http.Client{},
	}
}

func (h *lineHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	msg := h.formatMessage(changes)

	form := url.Values{}
	form.Set("message", msg)

	req, err := http.NewRequest(http.MethodPost, h.apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("line: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+h.token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("line: send request: %w", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("line: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func (h *lineHandler) formatMessage(changes []monitor.Change) string {
	var sb strings.Builder
	if h.prefix != "" {
		sb.WriteString(h.prefix)
		sb.WriteString(" ")
	}
	sb.WriteString(fmt.Sprintf("%d port change(s) detected:\n", len(changes)))
	for _, c := range changes {
		sb.WriteString(fmt.Sprintf("  %s\n", c.String()))
	}
	return sb.String()
}

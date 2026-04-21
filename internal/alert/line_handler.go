package alert

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// lineHandler sends alerts to LINE Notify.
type lineHandler struct {
	token   string
	apiURL  string
	prefix  string
	client  *http.Client
}

// NewLineHandler creates a new LINE Notify alert handler.
func NewLineHandler(token, apiURL, prefix string) Handler {
	if apiURL == "" {
		apiURL = "https://notify-api.line.me/api/notify"
	}
	if prefix == "" {
		prefix = "[portwatch]"
	}
	return &lineHandler{
		token:  token,
		apiURL: apiURL,
		prefix: prefix,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (h *lineHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	var sb strings.Builder
	sb.WriteString(h.prefix)
	sb.WriteString(fmt.Sprintf(" %d port change(s) detected:\n", len(changes)))
	for _, c := range changes {
		sb.WriteString("  " + c.String() + "\n")
	}

	form := url.Values{}
	form.Set("message", sb.String())

	req, err := http.NewRequest(http.MethodPost, h.apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("line: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+h.token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("line: send request: %w", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, resp.Body); resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("line: unexpected status %d", resp.StatusCode)
	}
	return nil
}

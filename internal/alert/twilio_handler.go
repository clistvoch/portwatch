package alert

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/deanrtaylor1/portwatch/internal/monitor"
)

type twilioHandler struct {
	accountSID string
	authToken  string
	from       string
	to         string
	baseURL    string
	client     *http.Client
}

// NewTwilioHandler returns a Handler that sends SMS alerts via the Twilio API.
func NewTwilioHandler(accountSID, authToken, from, to, baseURL string) Handler {
	return &twilioHandler{
		accountSID: accountSID,
		authToken:  authToken,
		from:       from,
		to:         to,
		baseURL:    strings.TrimRight(baseURL, "/"),
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (h *twilioHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	body := formatTwilioBody(changes)
	endpoint := fmt.Sprintf("%s/2010-04-01/Accounts/%s/Messages.json", h.baseURL, h.accountSID)

	form := url.Values{}
	form.Set("From", h.from)
	form.Set("To", h.to)
	form.Set("Body", body)

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("twilio: build request: %w", err)
	}
	req.SetBasicAuth(h.accountSID, h.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("twilio: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		var apiErr struct {
			Message string `json:"message"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&apiErr)
		return fmt.Errorf("twilio: unexpected status %d: %s", resp.StatusCode, apiErr.Message)
	}
	return nil
}

func formatTwilioBody(changes []monitor.Change) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("portwatch: %d port change(s) detected\n", len(changes)))
	for _, c := range changes {
		sb.WriteString(fmt.Sprintf("  %s\n", c.String()))
	}
	return sb.String()
}

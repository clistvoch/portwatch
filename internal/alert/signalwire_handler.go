package alert

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/user/portwatch/internal/monitor"
)

// signalWireHandler sends SMS alerts via the SignalWire REST API.
type signalWireHandler struct {
	projectID string
	apiToken  string
	from      string
	to        string
	baseURL   string
	client    *http.Client
}

// NewSignalWireHandler constructs a Handler that delivers SMS messages
// through SignalWire whenever port changes are detected.
func NewSignalWireHandler(projectID, apiToken, from, to, spaceURL string) Handler {
	return &signalWireHandler{
		projectID: projectID,
		apiToken:  apiToken,
		from:      from,
		to:        to,
		baseURL:   fmt.Sprintf("https://%s/api/laml/2010-04-01/Accounts/%s/Messages.json", spaceURL, projectID),
		client:    &http.Client{},
	}
}

func (h *signalWireHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	body := formatSignalWireBody(changes)
	return h.send(body)
}

func (h *signalWireHandler) send(body string) error {
	form := url.Values{}
	form.Set("From", h.from)
	form.Set("To", h.to)
	form.Set("Body", body)

	req, err := http.NewRequest(http.MethodPost, h.baseURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("signalwire: create request: %w", err)
	}
	req.SetBasicAuth(h.projectID, h.apiToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("signalwire: send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("signalwire: unexpected status %d: %s", resp.StatusCode, raw)
	}
	return nil
}

func formatSignalWireBody(changes []monitor.Change) string {
	type entry struct {
		Port  int    `json:"port"`
		Proto string `json:"proto"`
		Kind  string `json:"kind"`
	}
	entries := make([]entry, 0, len(changes))
	for _, c := range changes {
		entries = append(entries, entry{Port: c.Port, Proto: c.Proto, Kind: c.Kind.String()})
	}
	out, err := json.Marshal(entries)
	if err != nil {
		return fmt.Sprintf("portwatch: %d port change(s) detected", len(changes))
	}
	return fmt.Sprintf("portwatch alert: %s", out)
}

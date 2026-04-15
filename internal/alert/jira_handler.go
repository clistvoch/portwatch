package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/wiring/portwatch/internal/monitor"
)

// JiraHandler creates a Jira issue for each batch of port changes.
type JiraHandler struct {
	baseURL    string
	username   string
	apiToken   string
	projectKey string
	issueType  string
	priority   string
	client     *http.Client
}

func NewJiraHandler(baseURL, username, apiToken, projectKey, issueType, priority string) *JiraHandler {
	return &JiraHandler{
		baseURL:    strings.TrimRight(baseURL, "/"),
		username:   username,
		apiToken:   apiToken,
		projectKey: projectKey,
		issueType:  issueType,
		priority:   priority,
		client:     &http.Client{},
	}
}

func (h *JiraHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	summary := fmt.Sprintf("portwatch: %d port change(s) detected", len(changes))
	var sb strings.Builder
	for _, c := range changes {
		sb.WriteString(c.String())
		sb.WriteString("\n")
	}

	payload := map[string]any{
		"fields": map[string]any{
			"project":   map[string]string{"key": h.projectKey},
			"summary":   summary,
			"issuetype": map[string]string{"name": h.issueType},
			"priority":  map[string]string{"name": h.priority},
			"description": map[string]any{
				"type":    "doc",
				"version": 1,
				"content": []map[string]any{
					{"type": "paragraph", "content": []map[string]any{
						{"type": "text", "text": sb.String()},
					}},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("jira: marshal payload: %w", err)
	}

	url := h.baseURL + "/rest/api/3/issue"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("jira: build request: %w", err)
	}
	req.SetBasicAuth(h.username, h.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("jira: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("jira: unexpected status %d", resp.StatusCode)
	}
	return nil
}

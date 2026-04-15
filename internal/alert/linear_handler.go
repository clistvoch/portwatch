package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/username/portwatch/internal/monitor"
)

type LinearHandler struct {
	apiKey     string
	teamID     string
	projectID  string
	priority   int
	labelIDs   []string
	baseURL    string
	assigneeID string
	client     *http.Client
}

type linearIssueRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	TeamID      string   `json:"teamId"`
	ProjectID   string   `json:"projectId,omitempty"`
	Priority    int      `json:"priority"`
	LabelIDs    []string `json:"labelIds,omitempty"`
	AssigneeID  string   `json:"assigneeId,omitempty"`
}

func NewLinearHandler(apiKey, teamID, projectID, baseURL, assigneeID string, priority int, labelIDs []string) *LinearHandler {
	return &LinearHandler{
		apiKey:     apiKey,
		teamID:     teamID,
		projectID:  projectID,
		priority:   priority,
		labelIDs:   labelIDs,
		baseURL:    strings.TrimRight(baseURL, "/"),
		assigneeID: assigneeID,
		client:     &http.Client{},
	}
}

func (h *LinearHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	var sb strings.Builder
	for _, c := range changes {
		sb.WriteString(fmt.Sprintf("- %s\n", c.String()))
	}

	issue := linearIssueRequest{
		Title:       fmt.Sprintf("portwatch: %d port change(s) detected", len(changes)),
		Description: sb.String(),
		TeamID:      h.teamID,
		ProjectID:   h.projectID,
		Priority:    h.priority,
		LabelIDs:    h.labelIDs,
		AssigneeID:  h.assigneeID,
	}

	body, err := json.Marshal(issue)
	if err != nil {
		return fmt.Errorf("linear: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, h.baseURL+"/issues", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("linear: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", h.apiKey)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("linear: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("linear: unexpected status %d", resp.StatusCode)
	}
	return nil
}

package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

type matrixHandler struct {
	homeserver string
	token      string
	roomID     string
	msgType    string
	client     *http.Client
}

// NewMatrixHandler creates a Handler that sends port change alerts to a Matrix room.
func NewMatrixHandler(homeserver, token, roomID, msgType string) Handler {
	if msgType == "" {
		msgType = "m.text"
	}
	return &matrixHandler{
		homeserver: strings.TrimRight(homeserver, "/"),
		token:      token,
		roomID:     roomID,
		msgType:    msgType,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (h *matrixHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	var sb strings.Builder
	for _, c := range changes {
		sb.WriteString(c.String())
		sb.WriteString("\n")
	}

	txnID := fmt.Sprintf("portwatch-%d", time.Now().UnixNano())
	url := fmt.Sprintf("%s/_matrix/client/v3/rooms/%s/send/m.room.message/%s",
		h.homeserver, h.roomID, txnID)

	payload := map[string]string{
		"msgtype": h.msgType,
		"body":    sb.String(),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("matrix: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("matrix: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.token)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("matrix: send event: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("matrix: unexpected status %d", resp.StatusCode)
	}
	return nil
}

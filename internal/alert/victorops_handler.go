package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

type victorOpsPayload struct {
	MessageType       string `json:"message_type"`
	EntityID          string `json:"entity_id"`
	EntityDisplayName string `json:"entity_display_name"`
	StateMessage      string `json:"state_message"`
	Timestamp         int64  `json:"timestamp"`
}

type VictorOpsHandler struct {
	webhookURL  string
	routingKey  string
	messageType string
	client      *http.Client
}

func NewVictorOpsHandler(webhookURL, routingKey, messageType string) *VictorOpsHandler {
	return &VictorOpsHandler{
		webhookURL:  webhookURL,
		routingKey:  routingKey,
		messageType: messageType,
		client:      &http.Client{Timeout: 10 * time.Second},
	}
}

func (h *VictorOpsHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	body := fmt.Sprintf("portwatch detected %d port change(s)", len(changes))
	for _, c := range changes {
		body += fmt.Sprintf("\n  %s", c.String())
	}

	payload := victorOpsPayload{
		MessageType:       h.messageType,
		EntityID:          "portwatch.port-change",
		EntityDisplayName: fmt.Sprintf("Port changes detected (%d)", len(changes)),
		StateMessage:      body,
		Timestamp:         time.Now().Unix(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("victorops: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/%s", h.webhookURL, h.routingKey)
	resp, err := h.client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("victorops: send alert: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("victorops: unexpected status %d", resp.StatusCode)
	}
	return nil
}

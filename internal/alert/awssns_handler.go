package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/danvolchek/portwatch/internal/monitor"
)

// AWSSNSHandler publishes portwatch change events to an AWS SNS topic
// via the SNS HTTP publish endpoint (useful for testing/mocking).
// Production deployments should swap the HTTP call for the AWS SDK.
type AWSSNSHandler struct {
	endpoint  string
	topicARN  string
	subject   string
	accessKey string
	secretKey string
	client    *http.Client
}

type awsSNSPayload struct {
	TopicARN string `json:"TopicArn"`
	Subject  string `json:"Subject"`
	Message  string `json:"Message"`
}

func NewAWSSNSHandler(endpoint, topicARN, subject, accessKey, secretKey string) *AWSSNSHandler {
	if subject == "" {
		subject = "portwatch alert"
	}
	return &AWSSNSHandler{
		endpoint:  strings.TrimRight(endpoint, "/"),
		topicARN:  topicARN,
		subject:   subject,
		accessKey: accessKey,
		secretKey: secretKey,
		client:    &http.Client{},
	}
}

func (h *AWSSNSHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	var sb strings.Builder
	for _, c := range changes {
		sb.WriteString(c.String())
		sb.WriteByte('\n')
	}

	payload := awsSNSPayload{
		TopicARN: h.topicARN,
		Subject:  h.subject,
		Message:  sb.String(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("awssns: marshal payload: %w", err)
	}

	resp, err := h.client.Post(h.endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("awssns: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("awssns: unexpected status %d", resp.StatusCode)
	}
	return nil
}

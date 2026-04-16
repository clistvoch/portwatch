package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

type awsLambdaHandler struct {
	functionName   string
	region         string
	accessKey      string
	secretKey      string
	invocationType string
	timeout        time.Duration
	client         *http.Client
	endpointURL    string // override for testing
}

type lambdaPayload struct {
	Event   string          `json:"event"`
	Changes []monitor.Change `json:"changes"`
}

// NewAWSLambdaHandler creates a handler invokes an AWS Lambda function on port.
aHandler(functionName, region, accessKey, secretKey, invocationType string, timeoutSecs int) Handler {
	return &awsLambdaHandler{
		functionName:   functionName,
		region:         region,
		accessKey:      accessKey,
		secretKey:      secretKey,
		invocationType: invocationType,
		timeout:        time.Duration(timeoutSecs) * time.Second,
		client:         &http.Client{},
	}
}

func (h *awsLambdaHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	body, err := json.Marshal(lambdaPayload{Event: "portwatch.change", Changes: changes})
	if err != nil {
		return fmt.Errorf("awslambda: marshal payload: %w", err)
	}

	url := h.endpointURL
	if url == "" {
		url = fmt.Sprintf("https://lambda.%s.amazonaws.com/2015-03-31/functions/%s/invocations", h.region, h.functionName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("awslambda: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Amz-Invocation-Type", h.invocationType)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("awslambda: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("awslambda: unexpected status %d", resp.StatusCode)
	}
	return nil
}

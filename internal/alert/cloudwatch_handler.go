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

type CloudWatchHandler struct {
	region     string
	namespace  string
	metricName string
	accessKey  string
	secretKey  string
	client     *http.Client
	endpoint   string // overridable for testing
}

type cloudWatchPayload struct {
	Namespace  string             n	MetricName string             `json:"metric_name"`
	Value      "`
	Tim             `json:"timestamp"`
	Dimensions map[string]string  `json:"dimensions"`
	Changes    []cloudWatchChange `json:"changes"`
}

type cloudWatchChange struct {
	Port   int    `json:"port"`
	Status string `json:"status"`
}

func NewCloudWatchHandler(region, namespace, metricName, accessKey, secretKey string) *CloudWatchHandler {
	return &CloudWatchHandler{
		region:     region,
		namespace:  namespace,
		metricName: metricName,
		accessKey:  accessKey,
		secretKey:  secretKey,
		client:     &http.Client{Timeout: 10 * time.Second},
		endpoint:   fmt.Sprintf("https://monitoring.%s.amazonaws.com/", region),
	}
}

func (h *CloudWatchHandler) Handle(ctx context.Context, changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	cc := make([]cloudWatchChange, len(changes))
	for i, c := range changes {
		cc[i] = cloudWatchChange{Port: c.Port, Status: c.Status.String()}
	}

	payload := cloudWatchPayload{
		Namespace:  h.namespace,
		MetricName: h.metricName,
		Value:      float64(len(changes)),
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Dimensions: map[string]string{"Region": h.region},
		Changes:    cc,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("cloudwatch: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("cloudwatch: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Amz-Access-Key", h.accessKey)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("cloudwatch: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("cloudwatch: unexpected status %d", resp.StatusCode)
	}
	return nil
}

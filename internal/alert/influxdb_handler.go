package alert

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// InfluxDBHandler writes port change events to InfluxDB via the v2 write API
// using the line protocol format.
type InfluxDBHandler struct {
	url         string
	token       string
	org         string
	bucket      string
	measurement string
	client      *http.Client
}

// NewInfluxDBHandler creates an InfluxDBHandler with the supplied configuration.
func NewInfluxDBHandler(url, token, org, bucket, measurement string, timeoutSec int) *InfluxDBHandler {
	return &InfluxDBHandler{
		url:         url,
		token:       token,
		org:         org,
		bucket:      bucket,
		measurement: measurement,
		client:      &http.Client{Timeout: time.Duration(timeoutSec) * time.Second},
	}
}

// Handle sends each change as an InfluxDB line-protocol record.
func (h *InfluxDBHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	var buf bytes.Buffer
	ts := time.Now().UnixNano()
	for _, c := range changes {
		kind := "opened"
		if c.Type == monitor.Closed {
			kind = "closed"
		}
		fmt.Fprintf(&buf, "%s,port=%d,protocol=%s kind=%q %d\n",
			h.measurement, c.Port.Port, c.Port.Protocol, kind, ts)
	}

	endpoint := fmt.Sprintf("%s/api/v2/write?org=%s&bucket=%s&precision=ns",
		h.url, h.org, h.bucket)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, &buf)
	if err != nil {
		return fmt.Errorf("influxdb: build request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+h.token)
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("influxdb: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("influxdb: unexpected status %d", resp.StatusCode)
	}
	return nil
}

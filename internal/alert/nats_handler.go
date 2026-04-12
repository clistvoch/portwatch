package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/yourorg/portwatch/internal/monitor"
)

// NATSHandler publishes port change events to a NATS subject.
type NATSHandler struct {
	nc      *nats.Conn
	subject string
	logger  *log.Logger
}

type natsPayload struct {
	Timestamp string          `json:"timestamp"`
	Changes   []monitor.Change `json:"changes"`
}

// NewNATSHandler creates a NATSHandler connected to the given URL.
func NewNATSHandler(url, subject, username, password string, logger *log.Logger) (*NATSHandler, error) {
	opts := []nats.Option{
		nats.Timeout(5 * time.Second),
		nats.MaxReconnects(5),
	}
	if username != "" {
		opts = append(opts, nats.UserInfo(username, password))
	}

	nc, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("nats: connect to %s: %w", url, err)
	}

	if logger == nil {
		logger = log.Default()
	}

	return &NATSHandler{nc: nc, subject: subject, logger: logger}, nil
}

// Handle publishes changes to the configured NATS subject.
func (h *NATSHandler) Handle(_ context.Context, changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	payload := natsPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Changes:   changes,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("nats: marshal payload: %w", err)
	}

	if err := h.nc.Publish(h.subject, data); err != nil {
		return fmt.Errorf("nats: publish to %s: %w", h.subject, err)
	}

	h.logger.Printf("nats: published %d change(s) to subject %q", len(changes), h.subject)
	return nil
}

// Close drains and closes the NATS connection.
func (h *NATSHandler) Close() error {
	return h.nc.Drain()
}

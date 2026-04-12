package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/cskr/portwatch/internal/monitor"
)

// RedisPublisher abstracts the Redis Publish call for testing.
type RedisPublisher interface {
	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
	Close() error
}

type redisPayload struct {
	Timestamp string          `json:"timestamp"`
	Changes   []changePayload `json:"changes"`
}

type changePayload struct {
	Type string `json:"type"`
	Port int    `json:"port"`
	Proto string `json:"proto"`
}

// RedisHandler publishes port-change events to a Redis Pub/Sub topic.
type RedisHandler struct {
	client  RedisPublisher
	topic   string
	timeout time.Duration
}

// NewRedisHandler creates a RedisHandler using the provided publisher.
func NewRedisHandler(client RedisPublisher, topic string) *RedisHandler {
	if topic == "" {
		topic = "portwatch:changes"
	}
	return &RedisHandler{
		client:  client,
		topic:   topic,
		timeout: 5 * time.Second,
	}
}

// Handle publishes a JSON-encoded payload for each batch of changes.
func (h *RedisHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	payload := redisPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	for _, c := range changes {
		payload.Changes = append(payload.Changes, changePayload{
			Type:  c.String(),
			Port:  c.Port.Port,
			Proto: c.Port.Proto,
		})
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("redis handler: marshal payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	if cmd := h.client.Publish(ctx, h.topic, string(data)); cmd.Err() != nil {
		return fmt.Errorf("redis handler: publish: %w", cmd.Err())
	}
	return nil
}

// Close releases the underlying Redis connection.
func (h *RedisHandler) Close() error {
	return h.client.Close()
}

package alert

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// kafkaProducer is a minimal interface for sending Kafka messages,
// allowing easy substitution in tests.
type kafkaProducer interface {
	SendMessage(topic string, payload []byte) error
	Close() error
}

// KafkaHandler sends port change alerts to a Kafka topic.
type KafkaHandler struct {
	producer kafkaProducer
	topic    string
}

type kafkaPayload struct {
	Timestamp string           `json:"timestamp"`
	Changes   []kafkaChangeDTO `json:"changes"`
}

type kafkaChangeDTO struct {
	Port   int    `json:"port"`
	Action string `json:"action"`
}

// NewKafkaHandler creates a KafkaHandler using the provided producer and topic.
func NewKafkaHandler(producer kafkaProducer, topic string) *KafkaHandler {
	return &KafkaHandler{producer: producer, topic: topic}
}

// Handle publishes a JSON payload to the configured Kafka topic.
func (h *KafkaHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	dtos := make([]kafkaChangeDTO, 0, len(changes))
	for _, c := range changes {
		dtos = append(dtos, kafkaChangeDTO{
			Port:   c.Port,
			Action: strings.ToLower(c.Type.String()),
		})
	}

	payload := kafkaPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Changes:   dtos,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("kafka: marshal payload: %w", err)
	}

	if err := h.producer.SendMessage(h.topic, data); err != nil {
		return fmt.Errorf("kafka: send message: %w", err)
	}
	return nil
}

// Close releases underlying producer resources.
func (h *KafkaHandler) Close() error {
	return h.producer.Close()
}

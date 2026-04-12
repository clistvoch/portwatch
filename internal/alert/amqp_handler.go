package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/cjvnjde/portwatch/internal/monitor"
)

// AMQPHandler publishes port-change events to a RabbitMQ exchange.
type AMQPHandler struct {
	url        string
	exchange   string
	routingKey string
	durable    bool
}

// NewAMQPHandler returns an AMQPHandler configured from the provided options.
func NewAMQPHandler(url, exchange, routingKey string, durable bool) *AMQPHandler {
	return &AMQPHandler{
		url:        url,
		exchange:   exchange,
		routingKey: routingKey,
		durable:    durable,
	}
}

// Handle publishes each change as a JSON message to the configured exchange.
func (h *AMQPHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	conn, err := amqp.Dial(h.url)
	if err != nil {
		return fmt.Errorf("amqp: dial: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("amqp: open channel: %w", err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(h.exchange, "topic", h.durable, false, false, false, nil); err != nil {
		return fmt.Errorf("amqp: declare exchange: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, c := range changes {
		body, err := json.Marshal(c)
		if err != nil {
			return fmt.Errorf("amqp: marshal change: %w", err)
		}
		err = ch.PublishWithContext(ctx, h.exchange, h.routingKey, false, false,
			amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent,
				Body:         body,
			},
		)
		if err != nil {
			return fmt.Errorf("amqp: publish: %w", err)
		}
	}
	return nil
}

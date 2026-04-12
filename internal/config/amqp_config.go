package config

import "fmt"

// AMQPConfig holds configuration for AMQP (RabbitMQ) alert handler.
type AMQPConfig struct {
	Enabled    bool   `toml:"enabled"`
	URL        string `toml:"url"`
	Exchange   string `toml:"exchange"`
	RoutingKey string `toml:"routing_key"`
	Durable    bool   `toml:"durable"`
}

func defaultAMQPConfig() AMQPConfig {
	return AMQPConfig{
		Enabled:    false,
		URL:        "amqp://guest:guest@localhost:5672/",
		Exchange:   "portwatch",
		RoutingKey: "port.change",
		Durable:    true,
	}
}

func validateAMQP(c AMQPConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return &ValidationError{Field: "amqp.url", Message: "url is required"}
	}
	if c.Exchange == "" {
		return fmt.Errorf("amqp.exchange: exchange name is required")
	}
	if c.RoutingKey == "" {
		return fmt.Errorf("amqp.routing_key: routing key is required")
	}
	return nil
}

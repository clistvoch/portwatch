package config

import "fmt"

// KafkaConfig holds configuration for the Kafka alert handler.
type KafkaConfig struct {
	Enabled  bool     `toml:"enabled"`
	Brokers  []string `toml:"brokers"`
	Topic    string   `toml:"topic"`
	ClientID string   `toml:"client_id"`
}

func defaultKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Enabled:  false,
		Brokers:  []string{},
		Topic:    "portwatch-alerts",
		ClientID: "portwatch",
	}
}

func validateKafka(c KafkaConfig) error {
	if !c.Enabled {
		return nil
	}
	if len(c.Brokers) == 0 {
		return ValidationError{Field: "kafka.brokers", Msg: "at least one broker address is required"}
	}
	for _, b := range c.Brokers {
		if b == "" {
			return ValidationError{Field: "kafka.brokers", Msg: "broker address must not be empty"}
		}
	}
	if c.Topic == "" {
		return ValidationError{Field: "kafka.topic", Msg: "topic must not be empty"}
	}
	if c.ClientID == "" {
		return fmt.Errorf("kafka.client_id must not be empty")
	}
	return nil
}

package config

import "fmt"

// NATSConfig holds configuration for the NATS alert handler.
type NATSConfig struct {
	Enabled  bool   `toml:"enabled"`
	URL      string `toml:"url"`
	Subject  string `toml:"subject"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

func defaultNATSConfig() NATSConfig {
	return NATSConfig{
		Enabled: false,
		URL:     "nats://localhost:4222",
		Subject: "portwatch.changes",
	}
}

func validateNATS(c NATSConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return fmt.Errorf("%w: nats.url is required", ErrInvalidConfig)
	}
	if c.Subject == "" {
		return fmt.Errorf("%w: nats.subject must not be empty", ErrInvalidConfig)
	}
	return nil
}

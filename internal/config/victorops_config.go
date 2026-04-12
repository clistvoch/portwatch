package config

import "fmt"

type VictorOpsConfig struct {
	Enabled    bool   `toml:"enabled"`
	WebhookURL string `toml:"webhook_url"`
	RoutingKey string `toml:"routing_key"`
	MessageType string `toml:"message_type"`
}

func defaultVictorOpsConfig() VictorOpsConfig {
	return VictorOpsConfig{
		Enabled:     false,
		MessageType: "CRITICAL",
	}
}

func validateVictorOps(cfg VictorOpsConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.WebhookURL == "" {
		return fmt.Errorf("victorops: webhook_url is required")
	}
	if cfg.RoutingKey == "" {
		return fmt.Errorf("victorops: routing_key is required")
	}
	valid := map[string]bool{"CRITICAL": true, "WARNING": true, "INFO": true, "RECOVERY": true}
	if !valid[cfg.MessageType] {
		return fmt.Errorf("victorops: invalid message_type %q (must be CRITICAL, WARNING, INFO, or RECOVERY)", cfg.MessageType)
	}
	return nil
}

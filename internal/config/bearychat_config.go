package config

// BearyChat (now Lark/Feishu predecessor) webhook notification config.
type BearyChatConfig struct {
	Enabled    bool   `toml:"enabled"`
	WebhookURL string `toml:"webhook_url"`
	Channel    string `toml:"channel"`
	Timeout    int    `toml:"timeout_seconds"`
}

func defaultBearyChatConfig() BearyChatConfig {
	return BearyChatConfig{
		Enabled:    false,
		WebhookURL: "",
		Channel:    "#general",
		Timeout:    5,
	}
}

func validateBearyChat(c BearyChatConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.WebhookURL == "" {
		return &ValidationError{Field: "bearychat.webhook_url", Message: "webhook URL is required"}
	}
	if c.Timeout <= 0 {
		return &ValidationError{Field: "bearychat.timeout_seconds", Message: "timeout must be positive"}
	}
	return nil
}

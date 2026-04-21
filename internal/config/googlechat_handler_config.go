package config

// GoogleChatHandlerConfig holds resolved configuration for the Google Chat handler.
type GoogleChatHandlerConfig struct {
	WebhookURL string
	Enabled    bool
}

func defaultGoogleChatHandlerConfig() GoogleChatHandlerConfig {
	return GoogleChatHandlerConfig{
		Enabled: false,
	}
}

func validateGoogleChatHandler(c GoogleChatHandlerConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.WebhookURL == "" {
		return &ValidationError{Field: "googlechat.webhook_url", Message: "must not be empty"}
	}
	return nil
}

// GoogleChatHandlerConfigFromSettings builds a GoogleChatHandlerConfig from
// the generic settings map stored in the main Config.
func GoogleChatHandlerConfigFromSettings(s map[string]string) GoogleChatHandlerConfig {
	cfg := defaultGoogleChatHandlerConfig()
	if v, ok := s["webhook_url"]; ok && v != "" {
		cfg.WebhookURL = v
		cfg.Enabled = true
	}
	if v, ok := s["enabled"]; ok {
		cfg.Enabled = v == "true"
	}
	return cfg
}

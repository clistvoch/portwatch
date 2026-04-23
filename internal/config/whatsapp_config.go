package config

import "fmt"

type WhatsAppConfig struct {
	Enabled     bool   `toml:"enabled"`
	Token       string `toml:"token"`
	PhoneID     string `toml:"phone_id"`
	Recipient   string `toml:"recipient"`
	APIBase     string `toml:"api_base"`
	TemplateMsg bool   `toml:"template_message"`
}

func defaultWhatsAppConfig() WhatsAppConfig {
	return WhatsAppConfig{
		Enabled:     false,
		APIBase:     "https://graph.facebook.com/v18.0",
		TemplateMsg: false,
	}
}

func validateWhatsApp(c WhatsAppConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.Token == "" {
		return &ValidationError{Field: "whatsapp.token", Msg: "token is required"}
	}
	if c.PhoneID == "" {
		return &ValidationError{Field: "whatsapp.phone_id", Msg: "phone_id is required"}
	}
	if c.Recipient == "" {
		return &ValidationError{Field: "whatsapp.recipient", Msg: "recipient is required"}
	}
	if c.APIBase == "" {
		return fmt.Errorf("whatsapp.api_base must not be empty")
	}
	return nil
}

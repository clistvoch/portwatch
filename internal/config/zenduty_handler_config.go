package config

// ZendutyHandlerConfig holds runtime-resolved Zenduty alert settings.
type ZendutyHandlerConfig struct {
	APIKey    string
	ServiceID string
	AlertType string
	Title     string
	Enabled   bool
}

func defaultZendutyHandlerConfig() ZendutyHandlerConfig {
	return ZendutyHandlerConfig{
		AlertType: "acknowledged",
		Title:     "portwatch: port change detected",
		Enabled:   false,
	}
}

func validateZendutyHandler(c ZendutyHandlerConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.APIKey == "" {
		return &ValidationError{Field: "zenduty.api_key", Message: "api_key is required"}
	}
	if c.ServiceID == "" {
		return &ValidationError{Field: "zenduty.service_id", Message: "service_id is required"}
	}
	validTypes := map[string]bool{
		"acknowledged": true,
		"resolved":     true,
		"triggered":    true,
	}
	if !validTypes[c.AlertType] {
		return &ValidationError{Field: "zenduty.alert_type", Message: "must be one of: acknowledged, resolved, triggered"}
	}
	return nil
}

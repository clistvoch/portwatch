package config

// PagerDutyConfig holds settings for the PagerDuty alert handler.
type PagerDutyConfig struct {
	Enabled    bool   `toml:"enabled"`
	RoutingKey string `toml:"routing_key"`
}

// defaultPagerDutyConfig returns a PagerDutyConfig with safe defaults.
func defaultPagerDutyConfig() PagerDutyConfig {
	return PagerDutyConfig{
		Enabled:    false,
		RoutingKey: "",
	}
}

// validatePagerDuty returns an error if the PagerDuty config is invalid.
func validatePagerDuty(cfg PagerDutyConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.RoutingKey == "" {
		return &ValidationError{Field: "pagerduty.routing_key", Reason: "must not be empty when pagerduty is enabled"}
	}
	return nil
}

// IsReady reports whether the PagerDuty integration is enabled and fully
// configured. It can be used by alert handlers to guard against sending
// events when the integration is not set up.
func (p PagerDutyConfig) IsReady() bool {
	return p.Enabled && p.RoutingKey != ""
}

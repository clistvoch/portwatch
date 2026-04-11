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

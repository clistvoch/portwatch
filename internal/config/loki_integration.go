package config

// init registers Loki config defaults and validation into the global config
// pipeline so that Load and Validate are aware of the Loki section.
func init() {
	registerDefault(func(c *Config) {
		if c.Loki == (LokiConfig{}) {
			c.Loki = defaultLokiConfig()
		}
		if c.Loki.Labels == nil {
			c.Loki.Labels = defaultLokiConfig().Labels
		}
	})

	registerValidator(func(c Config) error {
		return validateLoki(c.Loki)
	})
}

package config

func init() {
	// Register New Relic defaults and validation into the config pipeline.
	registerDefaults = append(registerDefaults, func(c *Config) {
		c.NewRelic = defaultNewRelicConfig()
	})
	registerValidators = append(registerValidators, func(c Config) error {
		return validateNewRelic(c.NewRelic)
	})
}

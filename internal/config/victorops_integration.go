package config

func init() {
	// Register VictorOps defaults and validation into the config pipeline.
	registerDefaults = append(registerDefaults, func(c *Config) {
		c.VictorOps = defaultVictorOpsConfig()
	})
	registerValidators = append(registerValidators, func(c Config) error {
		return validateVictorOps(c.VictorOps)
	})
}

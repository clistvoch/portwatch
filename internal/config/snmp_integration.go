package config

func init() {
	// Register SNMP defaults and validation into the config pipeline.
	registerDefaults(func(c *Config) {
		c.SNMP = defaultSNMPConfig()
	})
	registerValidator(func(c Config) error {
		return validateSNMP(c.SNMP)
	})
}

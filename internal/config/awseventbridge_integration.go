package config

func init() {
	registerDefault("awseventbridge", func(c *Config) {
		if c.AWSEventBridge == (AWSEventBridgeConfig{}) {
			c.AWSEventBridge = defaultAWSEventBridgeConfig()
		}
	})
	registerValidator(func(c Config) error {
		return validateAWSEventBridge(c.AWSEventBridge)
	})
}

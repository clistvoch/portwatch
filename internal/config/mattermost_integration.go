package config

func init() {
	// Register Mattermost defaults and validation into the global config
	// lifecycle. Called automatically on package import.
	registerDefaults(func(c *Config) {
		c.Mattermost = defaultMattermostConfig()
	})
	registerValidator(func(c Config) error {
		return validateMattermost(c.Mattermost)
	})
}

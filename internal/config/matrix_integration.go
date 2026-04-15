package config

func init() {
	registerDefault(func(c *Config) {
		c.Matrix = defaultMatrixConfig()
	})
	registerValidator(func(c Config) error {
		return validateMatrix(c.Matrix)
	})
}

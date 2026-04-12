package config

func init() {
	// Register InfluxDB defaults and validation into the config pipeline.
	registerDefault(func(c *Config) {
		c.InfluxDB = defaultInfluxDBConfig()
	})
	registerValidator(func(c Config) error {
		return validateInfluxDB(c.InfluxDB)
	})
}

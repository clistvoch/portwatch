package config

// init registers Kafka config defaults and validation into the global
// Config pipeline so it is automatically applied on Load.
func init() {
	registerDefault(func(c *Config) {
		if c.Kafka == (KafkaConfig{}) {
			c.Kafka = defaultKafkaConfig()
		}
	})

	registerValidator(func(c Config) error {
		return validateKafka(c.Kafka)
	})
}

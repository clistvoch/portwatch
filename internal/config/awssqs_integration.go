package config

func init() {
	registerDefault("awssqs", func(c *Settings) {
		if c.AWSSQS == (AWSSQSConfig{}) {
			c.AWSSQS = defaultAWSSQSConfig()
		}
	})
	registerValidator("awssqs", func(c Settings) error {
		return validateAWSSQS(c.AWSSQS)
	})
}

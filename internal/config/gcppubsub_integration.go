package config

func init() {
	registerDefault("gcppubsub", defaultGCPPubSubConfig)
	registerValidator("gcppubsub", func(c *Config) error {
		return validateGCPPubSub(c.GCPPubSub)
	})
}

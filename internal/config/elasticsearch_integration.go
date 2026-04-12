package config

func init() {
	// Register Elasticsearch defaults and validation into the config pipeline.
	registerDefault(func(c *Config) {
		c.Elasticsearch = defaultElasticsearchConfig()
	})
	registerValidator(func(c Config) error {
		return validateElasticsearch(c.Elasticsearch)
	})
}

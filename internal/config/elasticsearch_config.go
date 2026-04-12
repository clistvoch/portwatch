package config

import "fmt"

// ElasticsearchConfig holds settings for the Elasticsearch alert handler.
type ElasticsearchConfig struct {
	Enabled   bool   `toml:"enabled"`
	URL       string `toml:"url"`
	Index     string `toml:"index"`
	Username  string `toml:"username"`
	Password  string `toml:"password"`
	TimeoutSec int   `toml:"timeout_sec"`
}

func defaultElasticsearchConfig() ElasticsearchConfig {
	return ElasticsearchConfig{
		Enabled:    false,
		URL:        "http://localhost:9200",
		Index:      "portwatch",
		TimeoutSec: 5,
	}
}

func validateElasticsearch(c ElasticsearchConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return &ValidationError{Field: "elasticsearch.url", Msg: "url is required"}
	}
	if c.Index == "" {
		return &ValidationError{Field: "elasticsearch.index", Msg: "index is required"}
	}
	if c.TimeoutSec <= 0 {
		return fmt.Errorf("elasticsearch.timeout_sec must be positive")
	}
	return nil
}

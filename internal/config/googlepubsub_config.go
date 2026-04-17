package config

import "fmt"

type GooglePubSubConfig struct {
	Enabled   bool   `toml:"enabled"`
	ProjectID string `toml:"project_id"`
	TopicID   string `toml:"topic_id"`
	CredsFile string `toml:"credentials_file"`
}

func defaultGooglePubSubConfig() GooglePubSubConfig {
	return GooglePubSubConfig{
		Enabled:   false,
		ProjectID: "",
		TopicID:   "",
		CredsFile: "",
	}
}

func validateGooglePubSub(c GooglePubSubConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.ProjectID == "" {
		return fmt.Errorf("googlepubsub: project_id is required")
	}
	if c.TopicID == "" {
		return fmt.Errorf("googlepubsub: topic_id is required")
	}
	return nil
}

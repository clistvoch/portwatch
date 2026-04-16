package config

import "fmt"

type GCPPubSubConfig struct {
	Enabled   bool   `toml:"enabled"`
	ProjectID string `toml:"project_id"`
	TopicID   string `toml:"topic_id"`
	CredFile  string `toml:"credentials_file"`
}

func defaultGCPPubSubConfig() GCPPubSubConfig {
	return GCPPubSubConfig{
		Enabled:   false,
		ProjectID: "",
		TopicID:   "portwatch-alerts",
		CredFile:  "",
	}
}

func validateGCPPubSub(c GCPPubSubConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.ProjectID == "" {
		return fmt.Errorf("gcppubsub: project_id is required")
	}
	if c.TopicID == "" {
		return fmt.Errorf("gcppubsub: topic_id is required")
	}
	return nil
}

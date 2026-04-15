package config

import "fmt"

type JiraConfig struct {
	Enabled    bool   `toml:"enabled"`
	BaseURL    string `toml:"base_url"`
	Username   string `toml:"username"`
	APIToken   string `toml:"api_token"`
	ProjectKey string `toml:"project_key"`
	IssueType  string `toml:"issue_type"`
	Priority   string `toml:"priority"`
}

func defaultJiraConfig() JiraConfig {
	return JiraConfig{
		Enabled:   false,
		IssueType: "Bug",
		Priority:  "High",
	}
}

func validateJira(cfg JiraConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.BaseURL == "" {
		return fmt.Errorf("jira: base_url is required")
	}
	if cfg.Username == "" {
		return fmt.Errorf("jira: username is required")
	}
	if cfg.APIToken == "" {
		return fmt.Errorf("jira: api_token is required")
	}
	if cfg.ProjectKey == "" {
		return fmt.Errorf("jira: project_key is required")
	}
	validPriorities := map[string]bool{"Highest": true, "High": true, "Medium": true, "Low": true, "Lowest": true}
	if !validPriorities[cfg.Priority] {
		return fmt.Errorf("jira: invalid priority %q, must be one of Highest, High, Medium, Low, Lowest", cfg.Priority)
	}
	return nil
}

package config

import (
	"os"
	"testing"
)

func writeJiraConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_JiraDefaults(t *testing.T) {
	cfg := defaultJiraConfig()
	if cfg.Enabled {
		t.Error("expected enabled=false by default")
	}
	if cfg.IssueType != "Bug" {
		t.Errorf("expected issue_type=Bug, got %s", cfg.IssueType)
	}
	if cfg.Priority != "High" {
		t.Errorf("expected priority=High, got %s", cfg.Priority)
	}
}

func TestLoad_JiraSection(t *testing.T) {
	path := writeJiraConfig(t, `
[jira]
enabled = true
base_url = "https://example.atlassian.net"
username = "user@example.com"
api_token = "tok123"
project_key = "OPS"
issue_type = "Task"
priority = "Medium"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Jira.Enabled {
		t.Error("expected enabled=true")
	}
	if cfg.Jira.BaseURL != "https://example.atlassian.net" {
		t.Errorf("unexpected base_url: %s", cfg.Jira.BaseURL)
	}
	if cfg.Jira.ProjectKey != "OPS" {
		t.Errorf("unexpected project_key: %s", cfg.Jira.ProjectKey)
	}
	if cfg.Jira.Priority != "Medium" {
		t.Errorf("unexpected priority: %s", cfg.Jira.Priority)
	}
}

func TestLoad_JiraMissingAPIKey(t *testing.T) {
	cfg := JiraConfig{
		Enabled:    true,
		BaseURL:    "https://example.atlassian.net",
		Username:   "user@example.com",
		ProjectKey: "OPS",
		Priority:   "High",
	}
	if err := validateJira(cfg); err == nil {
		t.Error("expected error for missing api_token")
	}
}

func TestLoad_JiraInvalidPriority(t *testing.T) {
	cfg := JiraConfig{
		Enabled:    true,
		BaseURL:    "https://example.atlassian.net",
		Username:   "user@example.com",
		APIToken:   "tok",
		ProjectKey: "OPS",
		Priority:   "Critical",
	}
	if err := validateJira(cfg); err == nil {
		t.Error("expected error for invalid priority")
	}
}

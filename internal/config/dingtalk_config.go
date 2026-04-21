package config

import "fmt"

type DingTalkConfig struct {
	Enabled    bool   `toml:"enabled"`
	WebhookURL string `toml:"webhook_url"`
	Secret     string `toml:"secret"`
	MsgType    string `toml:"msg_type"`
}

func defaultDingTalkConfig() DingTalkConfig {
	return DingTalkConfig{
		Enabled: false,
		MsgType: "text",
	}
}

func validateDingTalk(c DingTalkConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.WebhookURL == "" {
		return &ValidationError{Field: "dingtalk.webhook_url", Msg: "webhook URL is required"}
	}
	validTypes := map[string]bool{"text": true, "markdown": true}
	if !validTypes[c.MsgType] {
		return &ValidationError{
			Field: "dingtalk.msg_type",
			Msg:   fmt.Sprintf("invalid msg_type %q: must be text or markdown", c.MsgType),
		}
	}
	return nil
}

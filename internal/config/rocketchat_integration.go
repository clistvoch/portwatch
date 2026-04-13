package config

func init() {
	registerDefaulter(func(c *Config) {
		if c.RocketChat == nil {
			c.RocketChat = defaultRocketChatConfig()
		}
	})
	registerValidator(func(c *Config) error {
		if c.RocketChat != nil && c.RocketChat.Enabled {
			return validateRocketChat(c.RocketChat)
		}
		return nil
	})
}

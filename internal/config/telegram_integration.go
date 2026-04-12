package config

// Ensure TelegramConfig is wired into the top-level Config struct and
// that its defaults and validation are invoked by Load/Validate.
// This file registers the Telegram sub-config with the config lifecycle.

func init() {
	// Register default factory.
	subDefaults = append(subDefaults, func(c *Config) {
		c.Telegram = defaultTelegramConfig()
	})

	// Register validator.
	subValidators = append(subValidators, func(c Config) error {
		return validateTelegram(c.Telegram)
	})
}

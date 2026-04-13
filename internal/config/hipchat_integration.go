package config

func init() {
	registerDefault("hipchat", defaultHipChatConfig)
	registerValidator("hipchat", validateHipChat)
}

package config

func init() {
	registerDefaulter(func(s *Settings) {
		if s.XMPP == (XMPPConfig{}) {
			s.XMPP = defaultXMPPConfig()
		}
	})

	registerValidator(func(s Settings) error {
		return validateXMPP(s.XMPP)
	})
}

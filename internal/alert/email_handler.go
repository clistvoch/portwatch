package alert

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/user/portwatch/internal/monitor"
)

// EmailConfig holds SMTP configuration for the email handler.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

// emailHandler sends alert emails via SMTP.
type emailHandler struct {
	cfg  EmailConfig
	auth smtp.Auth
}

// NewEmailHandler creates a Handler that sends email alerts.
func NewEmailHandler(cfg EmailConfig) Handler {
	var auth smtp.Auth
	if cfg.Username != "" {
		auth = smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	}
	return &emailHandler{cfg: cfg, auth: auth}
}

func (e *emailHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Subject: portwatch: %d port change(s) detected\r\n", len(changes)))
	sb.WriteString(fmt.Sprintf("From: %s\r\n", e.cfg.From))
	sb.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(e.cfg.To, ", ")))
	sb.WriteString("\r\n")
	for _, c := range changes {
		sb.WriteString(c.String() + "\r\n")
	}

	addr := fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port)
	return smtp.SendMail(addr, e.auth, e.cfg.From, e.cfg.To, []byte(sb.String()))
}

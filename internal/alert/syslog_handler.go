package alert

import (
	"fmt"
	"log/syslog"
	"strings"

	"portwatch/internal/monitor"
)

// SyslogHandler sends alert notifications to the system syslog daemon.
type SyslogHandler struct {
	writer   *syslog.Writer
	priority syslog.Priority
	tag      string
}

// NewSyslogHandler creates a SyslogHandler that writes to syslog with the
// given priority and tag. Priority should be a syslog.Priority constant
// (e.g. syslog.LOG_WARNING | syslog.LOG_DAEMON).
func NewSyslogHandler(priority syslog.Priority, tag string) (*SyslogHandler, error) {
	if tag == "" {
		tag = "portwatch"
	}
	w, err := syslog.New(priority, tag)
	if err != nil {
		return nil, fmt.Errorf("syslog: open writer: %w", err)
	}
	return &SyslogHandler{writer: w, priority: priority, tag: tag}, nil
}

// Handle writes each change in the alert to syslog. It is a no-op when
// changes is empty.
func (h *SyslogHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	lines := make([]string, 0, len(changes))
	for _, c := range changes {
		lines = append(lines, c.String())
	}
	msg := fmt.Sprintf("portwatch: %d change(s) detected: %s",
		len(changes), strings.Join(lines, "; "))

	switch h.priority & 0x07 {
	case syslog.LOG_EMERG:
		return h.writer.Emerg(msg)
	case syslog.LOG_ALERT:
		return h.writer.Alert(msg)
	case syslog.LOG_CRIT:
		return h.writer.Crit(msg)
	case syslog.LOG_ERR:
		return h.writer.Err(msg)
	case syslog.LOG_WARNING:
		return h.writer.Warning(msg)
	case syslog.LOG_NOTICE:
		return h.writer.Notice(msg)
	case syslog.LOG_DEBUG:
		return h.writer.Debug(msg)
	default:
		return h.writer.Info(msg)
	}
}

// Close releases the underlying syslog connection.
func (h *SyslogHandler) Close() error {
	return h.writer.Close()
}

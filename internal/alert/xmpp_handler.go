package alert

import (
	"fmt"
	"net"
	"strings"

	"github.com/netwatch/portwatch/internal/monitor"
)

// XMPPHandler sends alert messages over XMPP (Jabber) using a plain TCP
// connection with minimal XMPP framing. For production use a full XMPP
// library is recommended; this implementation covers the basic alert path.
type XMPPHandler struct {
	server   string
	port     int
	username string
	password string
	to       string
	useTLS   bool
	dial     func(network, addr string) (net.Conn, error)
}

// NewXMPPHandler returns an XMPPHandler configured with the provided settings.
func NewXMPPHandler(server string, port int, username, password, to string, useTLS bool) *XMPPHandler {
	return &XMPPHandler{
		server:   server,
		port:     port,
		username: username,
		password: password,
		to:       to,
		useTLS:   useTLS,
		dial:     net.Dial,
	}
}

// Handle sends an XMPP message for each detected change.
func (h *XMPPHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	body := formatXMPPBody(changes)
	addr := fmt.Sprintf("%s:%d", h.server, h.port)
	conn, err := h.dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("xmpp: dial %s: %w", addr, err)
	}
	defer conn.Close()

	// Minimal XMPP message stanza — sufficient for alerting over a pre-auth
	// connection in test/mock scenarios.
	stanza := fmt.Sprintf(
		"<message to=%q type='chat'><body>%s</body></message>",
		h.to, body,
	)
	if _, err := fmt.Fprint(conn, stanza); err != nil {
		return fmt.Errorf("xmpp: send stanza: %w", err)
	}
	return nil
}

func formatXMPPBody(changes []monitor.Change) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("portwatch: %d port change(s) detected\n", len(changes)))
	for _, c := range changes {
		sb.WriteString(fmt.Sprintf("  %s\n", c.String()))
	}
	return sb.String()
}

package alert

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/wlynxg/portwatch/internal/monitor"
)

// IRCHandler sends port change alerts to an IRC channel.
type IRCHandler struct {
	server  string
	port    int
	nick    string
	channel string
	password string
	useTLS  bool
	dial    func(network, addr string) (net.Conn, error)
}

// NewIRCHandler creates an IRCHandler with the provided settings.
func NewIRCHandler(server string, port int, nick, channel, password string, useTLS bool) *IRCHandler {
	dialFn := func(network, addr string) (net.Conn, error) {
		if useTLS {
			return tls.DialWithDialer(
				&net.Dialer{Timeout: 10 * time.Second},
				network, addr,
				&tls.Config{ServerName: server},
			)
		}
		return net.DialTimeout(network, addr, 10*time.Second)
	}
	return &IRCHandler{
		server:  server,
		port:    port,
		nick:    nick,
		channel: channel,
		password: password,
		useTLS:  useTLS,
		dial:    dialFn,
	}
}

// Handle sends a PRIVMSG to the configured IRC channel for each change.
func (h *IRCHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	addr := fmt.Sprintf("%s:%d", h.server, h.port)
	conn, err := h.dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("irc: dial %s: %w", addr, err)
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(15 * time.Second)) //nolint:errcheck

	send := func(line string) {
		fmt.Fprintf(conn, "%s\r\n", line)
	}

	if h.password != "" {
		send("PASS " + h.password)
	}
	send("NICK " + h.nick)
	send(fmt.Sprintf("USER %s 0 * :portwatch", h.nick))
	send("JOIN " + h.channel)

	var sb strings.Builder
	for _, c := range changes {
		sb.WriteString(c.String())
		sb.WriteString(" ")
	}
	msg := strings.TrimSpace(sb.String())
	send(fmt.Sprintf("PRIVMSG %s :[portwatch] %s", h.channel, msg))
	send("QUIT")
	return nil
}

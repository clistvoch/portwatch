package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// PortInfo holds metadata about an open port.
type PortInfo struct {
	Port     int
	Protocol string
	Address  string
}

// String returns a human-readable representation of a PortInfo.
func (p PortInfo) String() string {
	return fmt.Sprintf("%s:%d (%s)", p.Address, p.Port, p.Protocol)
}

// Scanner scans for open ports on the local machine.
type Scanner struct {
	Host      string
	PortRange [2]int
	Protocols []string
}

// NewScanner creates a Scanner with sensible defaults.
func NewScanner(host string, startPort, endPort int) *Scanner {
	return &Scanner{
		Host:      host,
		PortRange: [2]int{startPort, endPort},
		Protocols: []string{"tcp", "udp"},
	}
}

// Scan checks each port in the range and returns those that are open.
func (s *Scanner) Scan() ([]PortInfo, error) {
	if s.PortRange[0] < 1 || s.PortRange[1] > 65535 || s.PortRange[0] > s.PortRange[1] {
		return nil, fmt.Errorf("invalid port range: %d-%d", s.PortRange[0], s.PortRange[1])
	}

	var open []PortInfo
	for _, proto := range s.Protocols {
		for port := s.PortRange[0]; port <= s.PortRange[1]; port++ {
			addr := net.JoinHostPort(s.Host, strconv.Itoa(port))
			conn, err := net.Dial(proto, addr)
			if err != nil {
				if isRefused(err) {
					continue
				}
				continue
			}
			conn.Close()
			open = append(open, PortInfo{
				Port:     port,
				Protocol: proto,
				Address:  s.Host,
			})
		}
	}
	return open, nil
}

// isRefused returns true if the error indicates a connection was actively refused.
func isRefused(err error) bool {
	return err != nil && strings.Contains(err.Error(), "refused")
}

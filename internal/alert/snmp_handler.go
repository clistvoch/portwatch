package alert

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// snmpTrap is a minimal UDP-based SNMP trap sender (SNMPv2c enterprise trap).
type snmpTrap struct {
	target    string
	community string
	logger    *log.Logger
}

// NewSNMPHandler returns a Handler that sends a UDP SNMP trap for each batch
// of changes. It uses a minimal hand-crafted SNMPv2c trap payload so that the
// package has no external dependencies.
func NewSNMPHandler(target string, port int, community string, logger *log.Logger) Handler {
	if logger == nil {
		logger = log.Default()
	}
	return &snmpTrap{
		target:    fmt.Sprintf("%s:%d", target, port),
		community: community,
		logger:    logger,
	}
}

func (s *snmpTrap) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	payload := s.buildTrap(changes)

	conn, err := net.DialTimeout("udp", s.target, 5*time.Second)
	if err != nil {
		return fmt.Errorf("snmp: dial %s: %w", s.target, err)
	}
	defer conn.Close()

	if _, err := conn.Write(payload); err != nil {
		return fmt.Errorf("snmp: write trap: %w", err)
	}

	s.logger.Printf("[snmp] sent trap to %s (%d change(s))", s.target, len(changes))
	return nil
}

// buildTrap builds a minimal BER-encoded SNMPv2c Trap-PDU.
// The trap carries a single OctetString varbind with a human-readable summary.
func (s *snmpTrap) buildTrap(changes []monitor.Change) []byte {
	summary := fmt.Sprintf("portwatch: %d port change(s) detected at %s",
		len(changes), time.Now().UTC().Format(time.RFC3339))

	// Encode summary as BER OctetString
	strBytes := []byte(summary)
	strTLV := append([]byte{0x04, byte(len(strBytes))}, strBytes...)

	// Minimal SNMPv2c message (version=1, community, PDU type 0xA7=Trap-PDU)
	version := []byte{0x02, 0x01, 0x01}
	comm := []byte(s.community)
	commTLV := append([]byte{0x04, byte(len(comm))}, comm...)

	// PDU contents: request-id, error-status, error-index, varbind-list
	pduContents := append([]byte{
		0x02, 0x01, 0x01, // request-id = 1
		0x02, 0x01, 0x00, // error-status = 0
		0x02, 0x01, 0x00, // error-index = 0
	}, wrapSequence(wrapSequence(strTLV))...)
	pduTLV := append([]byte{0xA7, byte(len(pduContents))}, pduContents...)

	msgContents := append(append(version, commTLV...), pduTLV...)
	return append([]byte{0x30, byte(len(msgContents))}, msgContents...)
}

func wrapSequence(data []byte) []byte {
	return append([]byte{0x30, byte(len(data))}, data...)
}

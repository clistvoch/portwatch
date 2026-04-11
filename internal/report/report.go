package report

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Snapshot represents a point-in-time report of port state and changes.
type Snapshot struct {
	Timestamp time.Time
	OpenPorts []monitor.PortInfo
	Changes   []monitor.Change
}

// Printer writes port snapshots to an output destination.
type Printer struct {
	out io.Writer
}

// NewPrinter returns a Printer that writes to w.
// If w is nil, os.Stdout is used.
func NewPrinter(w io.Writer) *Printer {
	if w == nil {
		w = os.Stdout
	}
	return &Printer{out: w}
}

// PrintSnapshot writes a formatted snapshot to the printer's output.
func (p *Printer) PrintSnapshot(s Snapshot) error {
	fmt.Fprintf(p.out, "=== Port Report [%s] ===\n", s.Timestamp.Format(time.RFC3339))

	if len(s.Changes) > 0 {
		fmt.Fprintln(p.out, "\nChanges detected:")
		for _, c := range s.Changes {
			fmt.Fprintf(p.out, "  %s\n", c)
		}
	}

	fmt.Fprintln(p.out, "\nCurrently open ports:")
	if len(s.OpenPorts) == 0 {
		fmt.Fprintln(p.out, "  (none)")
		return nil
	}

	tw := tabwriter.NewWriter(p.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "  PORT\tPROTO\tADDRESS")
	for _, pi := range s.OpenPorts {
		fmt.Fprintf(tw, "  %d\t%s\t%s\n", pi.Port, pi.Proto, pi.Address)
	}
	return tw.Flush()
}

// PrintChangesOnly writes only the change list, suitable for alert summaries.
func (p *Printer) PrintChangesOnly(changes []monitor.Change) {
	if len(changes) == 0 {
		fmt.Fprintln(p.out, "No changes detected.")
		return
	}
	for _, c := range changes {
		fmt.Fprintf(p.out, "%s\n", c)
	}
}

// PrintSummary writes a one-line summary of the snapshot, including the
// number of open ports and detected changes. Useful for compact log output.
func (p *Printer) PrintSummary(s Snapshot) {
	fmt.Fprintf(p.out, "[%s] open ports: %d, changes: %d\n",
		s.Timestamp.Format(time.RFC3339),
		len(s.OpenPorts),
		len(s.Changes),
	)
}

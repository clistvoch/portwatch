package alert

import (
	"io"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/report"
)

// ReportHandler is a Handler that writes formatted change reports
// using the report.Printer.
type ReportHandler struct {
	printer *report.Printer
}

// NewReportHandler returns a ReportHandler that writes to w.
func NewReportHandler(w io.Writer) *ReportHandler {
	return &ReportHandler{
		printer: report.NewPrinter(w),
	}
}

// Handle satisfies the Handler interface by printing a snapshot
// containing only the reported change.
func (h *ReportHandler) Handle(c monitor.Change) error {
	snapshot := report.Snapshot{
		Timestamp: time.Now(),
		Changes:   []monitor.Change{c},
		OpenPorts: nil,
	}
	return h.printer.PrintSnapshot(snapshot)
}

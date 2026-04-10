package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/user/portwatch/internal/monitor"
)

// Handler processes port change events.
type Handler interface {
	Handle(c monitor.Change)
}

// LogHandler writes change alerts to a writer (defaults to os.Stdout).
type LogHandler struct {
	Out    io.Writer
	Prefix string
}

// NewLogHandler returns a LogHandler writing to stdout.
func NewLogHandler(prefix string) *LogHandler {
	return &LogHandler{Out: os.Stdout, Prefix: prefix}
}

// Handle prints the change to the configured writer.
func (h *LogHandler) Handle(c monitor.Change) {
	prefix := h.Prefix
	if prefix == "" {
		prefix = "ALERT"
	}
	fmt.Fprintf(h.Out, "%s %s\n", prefix, c)
}

// Dispatcher fans change events out to multiple handlers.
type Dispatcher struct {
	handlers []Handler
}

// NewDispatcher creates a Dispatcher with the given handlers.
func NewDispatcher(handlers ...Handler) *Dispatcher {
	return &Dispatcher{handlers: handlers}
}

// Add registers an additional handler.
func (d *Dispatcher) Add(h Handler) {
	d.handlers = append(d.handlers, h)
}

// Dispatch sends the change to every registered handler.
func (d *Dispatcher) Dispatch(c monitor.Change) {
	for _, h := range d.handlers {
		h.Handle(c)
	}
}

// DrainChannel reads from ch until it is closed, dispatching each change.
func (d *Dispatcher) DrainChannel(ch <-chan monitor.Change) {
	for c := range ch {
		d.Dispatch(c)
	}
}

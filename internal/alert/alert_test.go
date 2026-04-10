package alert_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

func makeChange(port int, status string) monitor.Change {
	return monitor.Change{Port: port, Status: status, At: time.Now()}
}

func TestLogHandler_Handle(t *testing.T) {
	var buf bytes.Buffer
	h := &alert.LogHandler{Out: &buf, Prefix: "WARN"}
	h.Handle(makeChange(9090, "opened"))

	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected prefix WARN in output: %q", out)
	}
	if !strings.Contains(out, "9090") {
		t.Errorf("expected port 9090 in output: %q", out)
	}
	if !strings.Contains(out, "opened") {
		t.Errorf("expected status opened in output: %q", out)
	}
}

func TestLogHandler_DefaultPrefix(t *testing.T) {
	var buf bytes.Buffer
	h := &alert.LogHandler{Out: &buf}
	h.Handle(makeChange(80, "closed"))
	if !strings.Contains(buf.String(), "ALERT") {
		t.Errorf("expected default prefix ALERT")
	}
}

func TestDispatcher_MultipleHandlers(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	h1 := &alert.LogHandler{Out: &buf1, Prefix: "H1"}
	h2 := &alert.LogHandler{Out: &buf2, Prefix: "H2"}

	d := alert.NewDispatcher(h1, h2)
	d.Dispatch(makeChange(443, "opened"))

	if !strings.Contains(buf1.String(), "H1") {
		t.Error("handler 1 not called")
	}
	if !strings.Contains(buf2.String(), "H2") {
		t.Error("handler 2 not called")
	}
}

func TestDispatcher_DrainChannel(t *testing.T) {
	var buf bytes.Buffer
	h := &alert.LogHandler{Out: &buf, Prefix: "TEST"}
	d := alert.NewDispatcher(h)

	ch := make(chan monitor.Change, 3)
	ch <- makeChange(22, "opened")
	ch <- makeChange(22, "closed")
	ch <- makeChange(8080, "opened")
	close(ch)

	d.DrainChannel(ch)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d: %q", len(lines), buf.String())
	}
}

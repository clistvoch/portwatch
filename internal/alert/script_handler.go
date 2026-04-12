package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/wander/portwatch/internal/monitor"
)

// ScriptHandler runs an external script for each batch of changes.
type ScriptHandler struct {
	path    string
	timeout time.Duration
	logger  *log.Logger
}

// NewScriptHandler creates a ScriptHandler that executes the given script.
func NewScriptHandler(path string, timeoutSec int, logger *log.Logger) *ScriptHandler {
	if logger == nil {
		logger = log.Default()
	}
	return &ScriptHandler{
		path:    path,
		timeout: time.Duration(timeoutSec) * time.Second,
		logger:  logger,
	}
}

// Handle encodes changes as JSON and passes them to the script via stdin.
func (h *ScriptHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	payload, err := json.Marshal(changes)
	if err != nil {
		return fmt.Errorf("script handler: marshal: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, h.path)
	cmd.Stdin = bytes.NewReader(payload)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("script handler: exec %q: %w (output: %s)", h.path, err, string(out))
	}

	h.logger.Printf("[script] ran %q for %d change(s)", h.path, len(changes))
	return nil
}

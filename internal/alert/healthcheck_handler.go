package alert

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// HealthCheckHandler serves a simple HTTP liveness endpoint and tracks
// whether the last scan cycle completed without error.
type HealthCheckHandler struct {
	server  *http.Server
	healthy atomic.Bool
}

// NewHealthCheckHandler creates and starts a health check HTTP server on addr
// responding at path. Call Close to shut it down.
func NewHealthCheckHandler(addr, path string) (*HealthCheckHandler, error) {
	h := &HealthCheckHandler{}
	h.healthy.Store(true)

	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if h.healthy.Load() {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "ok")
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, "unhealthy")
		}
	})

	h.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() { _ = h.server.ListenAndServe() }()
	return h, nil
}

// SetHealthy updates the liveness state reported by the endpoint.
func (h *HealthCheckHandler) SetHealthy(ok bool) {
	h.healthy.Store(ok)
}

// Close gracefully shuts down the HTTP server.
func (h *HealthCheckHandler) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return h.server.Shutdown(ctx)
}

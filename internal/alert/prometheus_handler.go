package alert

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"portwatch/internal/monitor"
)

// PrometheusHandler exposes port-change metrics via a Prometheus HTTP endpoint.
type PrometheusHandler struct {
	server   *http.Server
	opened   prometheus.Counter
	closed   prometheus.Counter
	total    atomic.Int64
	registry *prometheus.Registry
}

// NewPrometheusHandler creates and starts a Prometheus metrics server.
func NewPrometheusHandler(address, path string) (*PrometheusHandler, error) {
	reg := prometheus.NewRegistry()

	opened := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "portwatch_ports_opened_total",
		Help: "Total number of ports detected as newly opened.",
	})
	closed := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "portwatch_ports_closed_total",
		Help: "Total number of ports detected as newly closed.",
	})

	if err := reg.Register(opened); err != nil {
		return nil, fmt.Errorf("register opened counter: %w", err)
	}
	if err := reg.Register(closed); err != nil {
		return nil, fmt.Errorf("register closed counter: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle(path, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	srv := &http.Server{Addr: address, Handler: mux}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("prometheus: server error: %v", err)
		}
	}()

	return &PrometheusHandler{
		server:   srv,
		opened:   opened,
		closed:   closed,
		registry: reg,
	}, nil
}

// Handle implements alert.Handler.
func (h *PrometheusHandler) Handle(changes []monitor.Change) error {
	for _, c := range changes {
		switch c.Type {
		case monitor.Opened:
			h.opened.Inc()
			h.total.Add(1)
		case monitor.Closed:
			h.closed.Inc()
			h.total.Add(1)
		}
	}
	return nil
}

// Total returns the cumulative number of port changes recorded by this handler.
func (h *PrometheusHandler) Total() int64 {
	return h.total.Load()
}

// Close shuts down the metrics HTTP server.
func (h *PrometheusHandler) Close() error {
	return h.server.Shutdown(context.Background())
}

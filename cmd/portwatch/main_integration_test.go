//go:build integration

package main

import (
	"context"
	"net"
	"testing"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/monitor"
	"portwatch/internal/scanner"
)

func listenTCP(t *testing.T, addr string) net.Listener {
	t.Helper()
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatalf("failed to listen on %s: %v", addr, err)
	}
	t.Cleanup(func() { ln.Close() })
	return ln
}

func TestIntegration_DetectsNewPort(t *testing.T) {
	sc, err := scanner.NewScanner(19000, 19010)
	if err != nil {
		t.Fatalf("NewScanner: %v", err)
	}

	changes := make(chan alert.Change, 10)
	handler := alert.HandlerFunc(func(c alert.Change) {
		changes <- c
	})
	dispatcher := alert.NewDispatcher()
	dispatcher.Register(handler)

	mon := monitor.New(sc, dispatcher)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// perform first scan to establish baseline
	if err := mon.Scan(ctx); err != nil {
		t.Fatalf("initial scan: %v", err)
	}

	// open a port after baseline
	_ = listenTCP(t, "127.0.0.1:19005")

	if err := mon.Scan(ctx); err != nil {
		t.Fatalf("second scan: %v", err)
	}

	select {
	case c := <-changes:
		if c.Port != 19005 {
			t.Errorf("expected port 19005, got %d", c.Port)
		}
		if c.Type != alert.Opened {
			t.Errorf("expected Opened, got %v", c.Type)
		}
	case <-time.After(2 * time.Second):
		t.Error("timed out waiting for change event")
	}
}

func TestIntegration_DetectsClosedPort(t *testing.T) {
	sc, err := scanner.NewScanner(19000, 19010)
	if err != nil {
		t.Fatalf("NewScanner: %v", err)
	}

	changes := make(chan alert.Change, 10)
	handler := alert.HandlerFunc(func(c alert.Change) {
		changes <- c
	})
	dispatcher := alert.NewDispatcher()
	dispatcher.Register(handler)

	mon := monitor.New(sc, dispatcher)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// open a port before baseline so it appears in the initial snapshot
	ln := listenTCP(t, "127.0.0.1:19007")

	if err := mon.Scan(ctx); err != nil {
		t.Fatalf("initial scan: %v", err)
	}

	// close the port so the next scan detects the removal
	ln.Close()

	if err := mon.Scan(ctx); err != nil {
		t.Fatalf("second scan: %v", err)
	}

	select {
	case c := <-changes:
		if c.Port != 19007 {
			t.Errorf("expected port 19007, got %d", c.Port)
		}
		if c.Type != alert.Closed {
			t.Errorf("expected Closed, got %v", c.Type)
		}
	case <-time.After(2 * time.Second):
		t.Error("timed out waiting for change event")
	}
}

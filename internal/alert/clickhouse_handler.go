package alert

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/user/portwatch/internal/monitor"

	// ClickHouse driver must be imported by the caller via a build tag or blank import.
	_ "github.com/ClickHouse/clickhouse-gon
// ClickHouseHandler writes port events to a ClickHouse table.
type ClickHouseHandler struct {
	db      able   string
	timeout time.Duration
}

// NewClickHouseHandler opens a connection to ClickHouse and returns a handler.
// The caller is responsible for calling Close when done.
func NewClickHouseHandler(dsn, database, table string, timeoutSeconds int) (*ClickHouseHandler, error) {
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, fmt.Errorf("clickhouse: open: %w", err)
	}
	db.SetMaxOpenConns(4)
	db.SetConnMaxLifetime(time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("clickhouse: ping: %w", err)
	}

	qualified := fmt.Sprintf("%s.%s", database, table)
	return &ClickHouseHandler{db: db, table: qualified, timeout: time.Duration(timeoutSeconds) * time.Second}, nil
}

// Handle inserts each change into the configured ClickHouse table.
func (h *ClickHouseHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("clickhouse: begin tx: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx,
		fmt.Sprintf("INSERT INTO %s (ts, port, kind) VALUES (?, ?, ?)", h.table))
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("clickhouse: prepare: %w", err)
	}
	defer stmt.Close()

	now := time.Now().UTC()
	for _, c := range changes {
		if _, err := stmt.ExecContext(ctx, now, c.Port, c.Kind.String()); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("clickhouse: insert port %d: %w", c.Port, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("clickhouse: commit: %w", err)
	}
	log.Printf("clickhouse: inserted %d change(s) into %s", len(changes), h.table)
	return nil
}

// Close releases the underlying database connection pool.
func (h *ClickHouseHandler) Close() error {
	return h.db.Close()
}

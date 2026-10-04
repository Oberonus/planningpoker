package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"

	"planningpoker/internal/domain/events"
)

const queryTimeout = 10 * time.Second

//go:embed migrations/001_initial.sql
var initialSchema string

// OpenPostgres connects to PostgreSQL and applies transactional schema migrations.
// Each service uses its own database and login, rather than sharing application tables.
func OpenPostgres(ctx context.Context, dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	config, err := pgx.ParseConfig(dsn)

	if err != nil {
		// Driver parse errors can contain credentials from the connection string.
		return nil, errors.New("invalid DATABASE_URL")
	}

	db := stdlib.OpenDB(*config)

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(3)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("connect to PostgreSQL: %w", err)
	}

	if err := migrate(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate PostgreSQL: %w", err)
	}

	return db, nil
}

func migrate(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() { _ = tx.Rollback() }()
	// Serialize startup migrations across processes sharing this application's database.
	if _, err := tx.ExecContext(ctx, "SELECT pg_advisory_xact_lock(7795346)"); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version integer PRIMARY KEY, applied_at timestamptz NOT NULL DEFAULT now())`); err != nil {
		return err
	}
	var version int

	if err := tx.QueryRowContext(ctx, "SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version); err != nil {
		return err
	}

	if version > 1 {
		return fmt.Errorf("unsupported database schema version %d", version)
	}

	if version == 0 {
		if _, err := tx.ExecContext(ctx, initialSchema); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES (1)"); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func publishEvents(bus events.EventBus, pending []events.DomainEvent) {
	for _, e := range pending {
		if err := bus.Publish(e); err != nil {
			// Match the existing at-most-once bus behavior. Durable delivery needs an outbox.
			logrus.Errorf("failed to publish a domain event: %v", err)
		}
	}
}

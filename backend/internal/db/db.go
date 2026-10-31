// Package db opens the Postgres connection pool and runs migrations.
package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // registers "pgx5://"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"

	"meditrack/migrations"
)

// Connect creates a connection pool and checks that the database answers.
func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	// Neon's free tier sleeps when idle, so give the first ping some time.
	pingCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}

// newMigrator builds a migrate instance that reads the embedded SQL files.
func newMigrator(databaseURL string) (*migrate.Migrate, error) {
	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return nil, err
	}
	// golang-migrate's pgx driver expects the "pgx5://" scheme.
	url := databaseURL
	for _, prefix := range []string{"postgresql://", "postgres://"} {
		if strings.HasPrefix(url, prefix) {
			url = "pgx5://" + strings.TrimPrefix(url, prefix)
			break
		}
	}
	return migrate.NewWithSourceInstance("iofs", src, url)
}

// MigrateUp applies all migrations that have not been applied yet.
func MigrateUp(databaseURL string) error {
	m, err := newMigrator(databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// MigrateDown rolls back the most recent `steps` migrations.
func MigrateDown(databaseURL string, steps int) error {
	m, err := newMigrator(databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()
	if err := m.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// MigrateVersion reports the current migration version.
func MigrateVersion(databaseURL string) (version uint, dirty bool, err error) {
	m, err := newMigrator(databaseURL)
	if err != nil {
		return 0, false, err
	}
	defer m.Close()
	version, dirty, err = m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, nil // nothing applied yet
	}
	return version, dirty, err
}

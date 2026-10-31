// Command api starts the MediTrack HTTP server.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"meditrack/internal/auth"
	"meditrack/internal/config"
	"meditrack/internal/db"
	"meditrack/internal/handlers"
	"meditrack/internal/models"
	"meditrack/internal/repository"
)

func main() {
	if err := run(); err != nil {
		slog.Error("startup failed", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// ctx is cancelled when you press Ctrl+C or the host sends SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if cfg.AutoMigrate {
		slog.Info("running migrations")
		if err := db.MigrateUp(cfg.DatabaseURL); err != nil {
			return err
		}
	}

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := seedAdmin(ctx, repository.NewUserRepo(pool), cfg); err != nil {
		return err
	}

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handlers.NewRouter(cfg, pool),
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Run the server in the background so we can wait for Ctrl+C below.
	errCh := make(chan error, 1)
	go func() {
		slog.Info("API listening", "port", cfg.Port)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		slog.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
	return nil
}

// seedAdmin creates the first admin from ADMIN_EMAIL / ADMIN_PASSWORD when
// there is no active admin yet. On later runs it does nothing.
func seedAdmin(ctx context.Context, users *repository.UserRepo, cfg *config.Config) error {
	count, err := users.CountActiveAdmins(ctx)
	if err != nil || count > 0 {
		return err
	}
	if cfg.AdminEmail == "" || len(cfg.AdminPassword) < 8 {
		slog.Warn("no admin user exists; set ADMIN_EMAIL and ADMIN_PASSWORD (min 8 chars) to create one")
		return nil
	}

	email := strings.ToLower(strings.TrimSpace(cfg.AdminEmail))
	hash, err := auth.HashPassword(cfg.AdminPassword)
	if err != nil {
		return err
	}
	admin := &models.User{
		Name: cfg.AdminName, Email: email, PasswordHash: hash,
		Role: models.RoleAdmin, IsActive: true,
	}
	if err := users.Create(ctx, admin); errors.Is(err, repository.ErrDuplicate) {
		slog.Warn("ADMIN_EMAIL already belongs to a non-admin or inactive user; not seeding", "email", email)
		return nil
	} else if err != nil {
		return err
	}
	slog.Info("seeded first admin user", "email", email)
	return nil
}

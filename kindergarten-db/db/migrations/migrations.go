// Package migrations
package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func RunMigrations(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect %w", err)
	}
	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("failed to run migrations %w", err)
	}
	slog.Info("migrations completed successfully")
	return nil
}

func MigrateDown(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect %w", err)
	}
	if err := goose.Down(db, "migrations"); err != nil {
		return fmt.Errorf("failed to rollback migrations %w", err)
	}
	slog.Info("migrations rollback completed successfully")
	return nil
}

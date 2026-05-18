package db

import (
	"database/sql"
	"embed"
	"fmt"
	"path"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

const migrationsPath = "migrations"

//go:embed migrations/*.sql
var migrations embed.FS

// Open opens the SQLite database at path and applies pending migrations.
func Open(dbPath string) (*sql.DB, error) {
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	if err := runMigrations(database); err != nil {
		_ = database.Close()
		return nil, err
	}

	return database, nil
}

func runMigrations(database *sql.DB) error {
	// Ensure the database can record which migrations have already run.
	if _, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	// Read embedded migration files from the package migrations directory.
	entries, err := migrations.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}

	// Apply migrations in filename order.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		// Migration entries must be SQL files, not directories.
		if entry.IsDir() {
			return fmt.Errorf("migration entry %s is a directory", entry.Name())
		}

		version := entry.Name()

		// Skip migrations already recorded in schema_migrations.
		applied, err := migrationApplied(database, version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		// Load the SQL for this migration from the embedded filesystem.
		migration, err := migrations.ReadFile(path.Join(migrationsPath, version))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", version, err)
		}

		// Run each migration and its tracking insert atomically.
		tx, err := database.Begin()
		if err != nil {
			return fmt.Errorf("begin migration %s: %w", version, err)
		}

		// Apply the migration SQL.
		if _, err := tx.Exec(string(migration)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("apply migration %s: %w", version, err)
		}

		// Record the migration version only after the SQL succeeds.
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %s: %w", version, err)
		}

		// Commit both the schema/data changes and the migration record.
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", version, err)
		}
	}

	return nil
}

func migrationApplied(database *sql.DB, version string) (bool, error) {
	var count int
	if err := database.QueryRow(
		"SELECT COUNT(1) FROM schema_migrations WHERE version = ?",
		version,
	).Scan(&count); err != nil {
		return false, fmt.Errorf("check migration %s: %w", version, err)
	}

	return count > 0, nil
}

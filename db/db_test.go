package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenCreatesDatabaseAndRunsMigrations(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "maas-ldap.db")

	database, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open(%q) error = %v", dbPath, err)
	}
	defer database.Close()

	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("stat database file: %v", err)
	}

	var migrationCount int
	if err := database.QueryRow("SELECT COUNT(1) FROM schema_migrations").Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("schema_migrations count = %d, want 1", migrationCount)
	}

	wantPasswords := map[string]string{
		"student1": "parola1",
		"tb171":    "parola1234",
	}

	for username, wantPassword := range wantPasswords {
		var gotPassword string
		err := database.QueryRow(
			"SELECT maas_password FROM maas_user_mappings WHERE username = ?",
			username,
		).Scan(&gotPassword)
		if err != nil {
			t.Fatalf("query seeded user %q: %v", username, err)
		}
		if gotPassword != wantPassword {
			t.Fatalf("password for %q = %q, want %q", username, gotPassword, wantPassword)
		}
	}
}

package maas

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestLoadUserMappings(t *testing.T) {
	database := openTestDatabase(t)
	insertMapping(t, database, "student1", "parola1")
	insertMapping(t, database, "tb171", "Muxime64")

	users, err := LoadUserMappings(database)
	if err != nil {
		t.Fatalf("LoadUserMappings() error = %v", err)
	}

	if users["student1"] != "parola1" {
		t.Fatalf("student1 password = %q, want parola1", users["student1"])
	}
	if users["tb171"] != "Muxime64" {
		t.Fatalf("tb171 password = %q, want Muxime64", users["tb171"])
	}
}

func TestLoadUserMappingsRejectsBlankUsername(t *testing.T) {
	database := openTestDatabase(t)
	insertMapping(t, database, "", "parola1")

	if _, err := LoadUserMappings(database); err == nil {
		t.Fatal("LoadUserMappings() error = nil, want error")
	}
}

func TestLoadUserMappingsRejectsBlankPassword(t *testing.T) {
	database := openTestDatabase(t)
	insertMapping(t, database, "student1", "")

	if _, err := LoadUserMappings(database); err == nil {
		t.Fatal("LoadUserMappings() error = nil, want error")
	}
}

func openTestDatabase(t *testing.T) *sql.DB {
	t.Helper()

	database, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite database: %v", err)
	}
	database.SetMaxOpenConns(1)
	t.Cleanup(func() {
		_ = database.Close()
	})

	_, err = database.Exec(`
		CREATE TABLE maas_user_mappings (
			username TEXT PRIMARY KEY,
			maas_password TEXT NOT NULL
		);
	`)
	if err != nil {
		t.Fatalf("create maas_user_mappings: %v", err)
	}

	return database
}

func insertMapping(t *testing.T, database *sql.DB, username string, password string) {
	t.Helper()

	_, err := database.Exec(
		"INSERT INTO maas_user_mappings (username, maas_password) VALUES (?, ?)",
		username,
		password,
	)
	if err != nil {
		t.Fatalf("insert mapping: %v", err)
	}
}

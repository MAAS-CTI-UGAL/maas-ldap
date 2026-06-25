package config

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadLogSettingsUsesStderrWhenLogPathUnset(t *testing.T) {
	t.Setenv(envLogPath, "")
	originalOutput := log.Writer()
	t.Cleanup(func() { log.SetOutput(originalOutput) })

	settings, err := loadLogSettings()
	if err != nil {
		t.Fatalf("loadLogSettings() returned error: %v", err)
	}
	if settings.FilePath != "" {
		t.Fatalf("FilePath = %q, want empty", settings.FilePath)
	}
	if settings.File != nil {
		t.Fatalf("File = %v, want nil", settings.File)
	}
	if log.Writer() != os.Stderr {
		t.Fatal("log output was not set to os.Stderr")
	}
}

func TestLoadLogSettingsOpensConfiguredFile(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "service.log")
	t.Setenv(envLogPath, logPath)
	originalOutput := log.Writer()
	t.Cleanup(func() { log.SetOutput(originalOutput) })

	settings, err := loadLogSettings()
	if err != nil {
		t.Fatalf("loadLogSettings() returned error: %v", err)
	}
	t.Cleanup(func() { _ = settings.File.Close() })

	if settings.FilePath != logPath {
		t.Fatalf("FilePath = %q, want %q", settings.FilePath, logPath)
	}
	if settings.File == nil {
		t.Fatal("File is nil")
	}

	log.Print("hello")
	if err := settings.File.Sync(); err != nil {
		t.Fatalf("Sync() returned error: %v", err)
	}

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile() returned error: %v", err)
	}
	if string(content) == "" {
		t.Fatal("log file is empty")
	}
}

func TestLoadLogSettingsReturnsOpenError(t *testing.T) {
	t.Setenv(envLogPath, filepath.Join(t.TempDir(), "missing", "service.log"))

	if _, err := loadLogSettings(); err == nil {
		t.Fatal("loadLogSettings() returned nil error")
	}
}

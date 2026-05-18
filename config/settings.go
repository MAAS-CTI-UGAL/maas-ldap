package config

import (
	"errors"
	"log"
	"os"
	"strings"
)

const (
	defaultListenAddress = ":42069"
	envDBPath            = "DB_PATH"
	envLogPath           = "LOG_PATH"
	envPort              = "PORT"
)

var errMissingDBPath = errors.New("application configuration is incomplete. Please set DB_PATH.")

// AppSettings contains application-wide routes and file paths.
type AppSettings struct {
	ListenAddress string
	LogFilePath   string
	DBPath        string
}

func loadAppSettings() AppSettings {
	listenAddress := loadListenAddress()

	logFilePath := envOrDefault(envLogPath, "")

	dbPath := os.Getenv(envDBPath)
	if dbPath == "" {
		log.Fatal(errMissingDBPath)
	}

	return AppSettings{
		ListenAddress: listenAddress,
		LogFilePath:   logFilePath,
		DBPath:        dbPath,
	}
}

func loadListenAddress() string {
	port := os.Getenv(envPort)
	if port == "" {
		return defaultListenAddress
	}
	if port[0] == ':' || strings.Contains(port, ":") {
		return port
	}
	return ":" + port
}

// envOrDefault lets operators override selected defaults without changing code.
func envOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

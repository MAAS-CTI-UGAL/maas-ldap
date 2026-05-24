package config

import (
	"errors"
	"log"
	"strings"
)

const (
	defaultListenAddress = ":9090"
	envDBPath            = "DB_PATH"
	envLogPath           = "LOG_PATH"
	envPort              = "PORT"
)

var errMissingDBPath = errors.New("application configuration is incomplete. Please set DB_PATH to the SQLite database file path.")

// AppSettings contains application-wide routes and file paths.
type AppSettings struct {
	ListenAddress string
	LogFilePath   string
	DBPath        string
}

func loadAppSettings() AppSettings {
	listenAddress := loadListenAddress()

	// If envLogPath is not set, Configure will log to stderr only.
	logFilePath := envOrDefault(envLogPath, "")

	dbPath, err := loadDBPath()
	if err != nil {
		log.Fatal(err)
	}

	return AppSettings{
		ListenAddress: listenAddress,
		LogFilePath:   logFilePath,
		DBPath:        dbPath,
	}
}

func loadDBPath() (string, error) {
	dbPath := getEnv(envDBPath)
	if dbPath == "" {
		return "", errMissingDBPath
	}
	return dbPath, nil
}

func loadListenAddress() string {
	port := envOrDefault(envPort, defaultListenAddress)
	if port[0] == ':' || strings.Contains(port, ":") {
		return port
	}
	return ":" + port
}

package config

import (
	"os"
	"strings"
)

func getEnv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

// envOrDefault lets operators override selected defaults without changing code.
func envOrDefault(key, defaultValue string) string {
	value := getEnv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

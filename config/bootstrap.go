package config

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
)

const (
	defaultListenAddress = ":9090"
	envPort              = "PORT"
)

// AppConfig contains the runtime objects needed by the HTTP handlers.
type AppConfig struct {
	ListenAddress string
	LDAP          LDAPConfig
	Log           LogSettings
}

// Bootstrap loads and validates environment-driven startup configuration.
func Bootstrap() AppConfig {
	// Startup config is fatal because the service cannot authenticate or proxy
	// without these values.
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	listenAddress := loadListenAddress()

	logSettings, err := loadLogSettings()
	if err != nil {
		log.Fatal(err)
	}

	ldapConfig := loadLDAPConfig()

	return AppConfig{
		ListenAddress: listenAddress,
		Log:           logSettings,
		LDAP:          ldapConfig,
	}
}

func loadListenAddress() string {
	port := envOrDefault(envPort, defaultListenAddress)
	if port[0] == ':' || strings.Contains(port, ":") {
		return port
	}
	return ":" + port
}

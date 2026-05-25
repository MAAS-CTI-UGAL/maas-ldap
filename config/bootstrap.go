package config

import (
	"errors"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

var (
	errMissingLDAPURL       = errors.New("LDAP configuration is incomplete. Please set LDAP_URL.")
	errMissingLDAPUPNSuffix = errors.New("LDAP configuration is incomplete. Please set LDAP_UPN_SUFFIX.")
	errMissingLDAPBaseDN    = errors.New("LDAP configuration is incomplete. Please set LDAP_BASE_DN.")
)

// AppConfig contains the runtime objects needed by the HTTP handlers.
type AppConfig struct {
	Settings   AppSettings
	LDAP       LDAPConfig
	HTTPClient *http.Client
}

// LDAPConfig contains the environment-driven LDAP connection settings.
type LDAPConfig struct {
	URL        string
	UPN_SUFFIX string
	BASE_DN    string
}

// Bootstrap loads and validates environment-driven startup configuration.
func Bootstrap() AppConfig {
	// Startup config is fatal because the service cannot authenticate or proxy
	// without these values.
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	appSettings := loadAppSettings()
	ldapConfig := loadLDAPConfig()

	return AppConfig{
		Settings:   appSettings,
		LDAP:       ldapConfig,
		HTTPClient: http.DefaultClient,
	}
}

// loadLDAPConfig loads LDAP configuration from environment variables.
func loadLDAPConfig() LDAPConfig {
	ldapURL := getEnv("LDAP_URL")
	ldapUPNSuffix := getEnv("LDAP_UPN_SUFFIX")
	ldapBASEDN := getEnv("LDAP_BASE_DN")

	if ldapURL == "" {
		log.Fatal(errMissingLDAPURL)
	}

	if ldapUPNSuffix == "" {
		log.Fatal(errMissingLDAPUPNSuffix)
	}

	if ldapBASEDN == "" {
		log.Fatal(errMissingLDAPBaseDN)
	}

	return LDAPConfig{
		URL:        ldapURL,
		UPN_SUFFIX: ldapUPNSuffix,
		BASE_DN:    ldapBASEDN,
	}
}

package config

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"maas-ldap/backends/maas"
	"maas-ldap/db"

	"github.com/joho/godotenv"
)

var (
	errMissingLDAPURL          = errors.New("LDAP configuration is incomplete. Please set LDAP_URL.")
	errMissingLDAPUPNSuffix    = errors.New("LDAP configuration is incomplete. Please set LDAP_UPN_SUFFIX.")
	errMissingLDAPBaseDN       = errors.New("LDAP configuration is incomplete. Please set LDAP_BASE_DN.")
	errMissingLDAPAllowedGroup = errors.New("LDAP configuration is incomplete. Please set LDAP_ALLOWED_GROUP.")
	errBackendURLMissingHost   = errors.New("backend URL must include scheme and host")
)

// AppConfig contains the runtime objects needed by the HTTP handlers.
type AppConfig struct {
	Settings   AppSettings
	LDAP       LDAPConfig
	MAAS       BackendConfig
	Users      map[string]UserMapping
	HTTPClient *http.Client
}

// LDAPConfig contains the environment-driven LDAP connection settings.
type LDAPConfig struct {
	URL           string
	UPN_SUFFIX    string
	BASE_DN       string
	ALLOWED_GROUP string
}

// BackendConfig contains a backend base URL and its derived endpoint URLs.
type BackendConfig struct {
	BaseURL string
	URLs    map[string]url.URL
}

// UserMapping contains target app credentials mapped to an LDAP username.
type UserMapping struct {
	Password string
}

// Bootstrap loads and validates all configuration required to start the app.
func Bootstrap() AppConfig {
	// Startup config is fatal because the service cannot authenticate or proxy
	// without these values.
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	appSettings := loadAppSettings()
	ldapConfig := loadLDAPConfig()
	maasConfig := loadBackendConfig("MAAS_URL", maas.EndpointPaths)

	users, err := loadUsersFromDB(appSettings.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	return AppConfig{
		Settings:   appSettings,
		LDAP:       ldapConfig,
		MAAS:       maasConfig,
		Users:      users,
		HTTPClient: http.DefaultClient,
	}
}

// loadLDAPConfig loads LDAP configuration from environment variables.
func loadLDAPConfig() LDAPConfig {
	ldapURL := getEnv("LDAP_URL")
	ldapUPNSuffix := getEnv("LDAP_UPN_SUFFIX")
	ldapBASEDN := getEnv("LDAP_BASE_DN")
	ldapALLOWEDGROUP := getEnv("LDAP_ALLOWED_GROUP")

	if ldapURL == "" {
		log.Fatal(errMissingLDAPURL)
	}

	if ldapUPNSuffix == "" {
		log.Fatal(errMissingLDAPUPNSuffix)
	}

	if ldapBASEDN == "" {
		log.Fatal(errMissingLDAPBaseDN)
	}

	if ldapALLOWEDGROUP == "" {
		log.Fatal(errMissingLDAPAllowedGroup)
	}

	return LDAPConfig{
		URL:           ldapURL,
		UPN_SUFFIX:    ldapUPNSuffix,
		BASE_DN:       ldapBASEDN,
		ALLOWED_GROUP: ldapALLOWEDGROUP,
	}
}

func loadBackendConfig(baseURLKey string, endpointPaths map[string]string) BackendConfig {
	baseURL := getEnv(baseURLKey)
	if baseURL == "" {
		log.Fatalf("backend configuration is incomplete. Please set %s.", baseURLKey)
	}

	urls := map[string]url.URL{}
	for endpointKey, endpointPath := range endpointPaths {
		endpointURL, err := buildBackendURL(baseURL, endpointPath)
		if err != nil {
			if errors.Is(err, errBackendURLMissingHost) {
				log.Fatalf("%s must include scheme and host", baseURLKey)
			}
			log.Fatalf("%s is invalid: %v", baseURLKey, err)
		}
		urls[endpointKey] = endpointURL
	}

	// Store endpoint URLs by key so handlers do not rebuild target URLs per request.
	return BackendConfig{
		BaseURL: baseURL,
		URLs:    urls,
	}
}

// loadUsersFromDB loads LDAP-username to target-password mappings from SQLite.
func loadUsersFromDB(path string) (map[string]UserMapping, error) {
	database, err := db.Open(path)
	if err != nil {
		return nil, err
	}
	defer database.Close()

	passwords, err := maas.LoadUserMappings(database)
	if err != nil {
		return nil, err
	}
	users := map[string]UserMapping{}
	for username, password := range passwords {
		users[username] = UserMapping{
			Password: password,
		}
	}

	return users, nil
}

// buildBackendURL combines a configured backend origin with one backend route.
func buildBackendURL(baseURL string, path string) (url.URL, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return url.URL{}, err
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return url.URL{}, errBackendURLMissingHost
	}
	// Preserve any backend base path while avoiding duplicate slashes.
	parsedURL.Path = strings.TrimRight(parsedURL.Path, "/") + path
	return *parsedURL, nil
}

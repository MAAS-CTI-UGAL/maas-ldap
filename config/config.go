package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"maas-ldap/backends/maas"

	"github.com/joho/godotenv"
)

var (
	errMissingLDAPURL          = errors.New("LDAP configuration is incomplete. Please set LDAP_URL.")
	errMissingLDAPUPNSuffix    = errors.New("LDAP configuration is incomplete. Please set LDAP_UPN_SUFFIX.")
	errMissingLDAPBaseDN       = errors.New("LDAP configuration is incomplete. Please set LDAP_BASE_DN.")
	errMissingLDAPAllowedGroup = errors.New("LDAP configuration is incomplete. Please set LDAP_ALLOWED_GROUP.")
	errBackendURLMissingHost   = errors.New("backend URL must include scheme and host")
	errEmptyUsernameMapping    = errors.New("users.json contains an empty username")
)

// AppConfig contains the runtime objects needed by the HTTP handlers.
type AppConfig struct {
	App        AppSettings
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
	Password string `json:"maas_password"`
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

	users, err := loadUsers(appSettings.UsersFile)
	if err != nil {
		log.Fatal(err)
	}

	return AppConfig{
		App:        appSettings,
		LDAP:       ldapConfig,
		MAAS:       maasConfig,
		Users:      users,
		HTTPClient: http.DefaultClient,
	}
}

// loadLDAPConfig loads LDAP configuration from environment variables.
func loadLDAPConfig() LDAPConfig {
	ldapURL := os.Getenv("LDAP_URL")
	ldapUPNSuffix := os.Getenv("LDAP_UPN_SUFFIX")
	ldapBASEDN := os.Getenv("LDAP_BASE_DN")
	ldapALLOWEDGROUP := os.Getenv("LDAP_ALLOWED_GROUP")

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
	baseURL := os.Getenv(baseURLKey)
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

// loadUsers decodes the LDAP-username to target-password mapping file.
func loadUsers(path string) (map[string]UserMapping, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open users file: %w", err)
	}
	defer file.Close()

	users := map[string]UserMapping{}
	if err := json.NewDecoder(file).Decode(&users); err != nil {
		return nil, fmt.Errorf("decode users file: %w", err)
	}

	for username, user := range users {
		if strings.TrimSpace(username) == "" {
			return nil, errEmptyUsernameMapping
		}
		if strings.TrimSpace(user.Password) == "" {
			return nil, fmt.Errorf("users.json mapping for %q has an empty password", username)
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

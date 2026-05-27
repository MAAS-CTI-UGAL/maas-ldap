package backends

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"maas-ldap/config"
)

// LoginHandlerFactory creates one backend login route handler.
type LoginHandlerFactory func(config.AppConfig, url.URL, string) http.HandlerFunc

// Definition contains the static configuration for one known backend.
type Definition struct {
	Name            string
	BaseURLEnv      string
	AllowedGroupEnv string
	LoginPath       string
	NewLoginHandler LoginHandlerFactory
}

// Config contains a loaded backend definition and its environment values.
type Config struct {
	Definition
	BaseURL      string
	Target       url.URL
	AllowedGroup string
}

// LoadConfig loads and validates one backend definition.
func LoadConfig(definition Definition) (Config, error) {
	if definition.NewLoginHandler == nil {
		return Config{}, fmt.Errorf("backend %q login handler is not configured", definition.Name)
	}

	baseURL, target, err := loadBackendTarget(definition.BaseURLEnv, definition.LoginPath)
	if err != nil {
		return Config{}, err
	}

	allowedGroup, err := LoadAllowedGroup(definition.AllowedGroupEnv)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Definition:   definition,
		BaseURL:      baseURL,
		Target:       target,
		AllowedGroup: allowedGroup,
	}, nil
}

func loadBackendTarget(baseURLEnvKey string, path string) (string, url.URL, error) {
	baseURL := strings.TrimSpace(os.Getenv(baseURLEnvKey))
	if baseURL == "" {
		return "", url.URL{}, fmt.Errorf("backend configuration is incomplete. Please set %s.", baseURLEnvKey)
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", url.URL{}, fmt.Errorf("%s is invalid: %w", baseURLEnvKey, err)
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", url.URL{}, fmt.Errorf("%s must include scheme and host", baseURLEnvKey)
	}

	return baseURL, buildBackendURL(*parsedURL, path), nil
}

func buildBackendURL(baseURL url.URL, path string) url.URL {
	// Preserve any backend base path while avoiding duplicate slashes.
	baseURL.Path = strings.TrimRight(baseURL.Path, "/") + path
	return baseURL
}

// LoadAllowedGroup loads the LDAP group allowed to access one backend.
func LoadAllowedGroup(allowedGroupEnvKey string) (string, error) {
	allowedGroup := strings.TrimSpace(os.Getenv(allowedGroupEnvKey))
	if allowedGroup == "" {
		return "", fmt.Errorf("backend configuration is incomplete. Please set %s.", allowedGroupEnvKey)
	}
	return allowedGroup, nil
}

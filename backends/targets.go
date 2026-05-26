package backends

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// Targets maps backend target keys to their derived URLs.
type Targets map[string]url.URL

// Backend contains shared backend configuration.
type Backend struct {
	BaseURL      string
	Targets      Targets
	AllowedGroup string
}

// LoadBackendConfig loads and validates shared backend configuration.
func LoadBackendConfig(baseURLEnvKey string, allowedGroupEnvKey string, paths map[string]string) (Backend, error) {
	baseURL, targets, err := LoadBackendTargets(baseURLEnvKey, paths)
	if err != nil {
		return Backend{}, err
	}

	allowedGroup, err := LoadAllowedGroup(allowedGroupEnvKey)
	if err != nil {
		return Backend{}, err
	}

	return Backend{
		BaseURL:      baseURL,
		Targets:      targets,
		AllowedGroup: allowedGroup,
	}, nil
}

// LoadBackendTargets loads and validates target URLs for one backend.
func LoadBackendTargets(baseURLEnvKey string, paths map[string]string) (string, Targets, error) {
	baseURL := strings.TrimSpace(os.Getenv(baseURLEnvKey))
	if baseURL == "" {
		return "", nil, fmt.Errorf("backend configuration is incomplete. Please set %s.", baseURLEnvKey)
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", nil, fmt.Errorf("%s is invalid: %w", baseURLEnvKey, err)
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", nil, fmt.Errorf("%s must include scheme and host", baseURLEnvKey)
	}

	targets := Targets{}
	for targetKey, path := range paths {
		targets[targetKey] = buildBackendURL(*parsedURL, path)
	}

	return baseURL, targets, nil
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

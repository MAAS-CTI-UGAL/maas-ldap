package backends

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// BackendTargets contains a backend base URL and its derived target URLs.
type BackendTargets struct {
	BaseURL string
	Targets map[string]url.URL
}

// LoadBackendTargets loads and validates target URLs for one backend.
func LoadBackendTargets(baseURLEnvKey string, paths map[string]string) (BackendTargets, error) {
	baseURL := strings.TrimSpace(os.Getenv(baseURLEnvKey))
	if baseURL == "" {
		return BackendTargets{}, fmt.Errorf("backend configuration is incomplete. Please set %s.", baseURLEnvKey)
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return BackendTargets{}, fmt.Errorf("%s is invalid: %w", baseURLEnvKey, err)
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return BackendTargets{}, fmt.Errorf("%s must include scheme and host", baseURLEnvKey)
	}

	targets := map[string]url.URL{}
	for targetKey, path := range paths {
		targets[targetKey] = buildBackendURL(*parsedURL, path)
	}

	return BackendTargets{
		BaseURL: baseURL,
		Targets: targets,
	}, nil
}

func buildBackendURL(baseURL url.URL, path string) url.URL {
	// Preserve any backend base path while avoiding duplicate slashes.
	baseURL.Path = strings.TrimRight(baseURL.Path, "/") + path
	return baseURL
}

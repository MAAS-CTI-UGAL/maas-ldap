package backends

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

var errBackendURLMissingHost = errors.New("backend URL must include scheme and host")

// BackendTargets contains a backend base URL and its derived target URLs.
type BackendTargets struct {
	BaseURL string
	Targets map[string]url.URL
}

// LoadBackendTargets loads and validates target URLs for one backend.
func LoadBackendTargets(envKey string, paths map[string]string) (BackendTargets, error) {
	baseURL := strings.TrimSpace(os.Getenv(envKey))
	if baseURL == "" {
		return BackendTargets{}, fmt.Errorf("backend configuration is incomplete. Please set %s.", envKey)
	}

	targets := map[string]url.URL{}
	for targetKey, path := range paths {
		targetURL, err := BuildBackendURL(baseURL, path)
		if err != nil {
			if errors.Is(err, errBackendURLMissingHost) {
				return BackendTargets{}, fmt.Errorf("%s must include scheme and host", envKey)
			}
			return BackendTargets{}, fmt.Errorf("%s is invalid: %w", envKey, err)
		}
		targets[targetKey] = targetURL
	}

	return BackendTargets{
		BaseURL: baseURL,
		Targets: targets,
	}, nil
}

// BuildBackendURL combines a configured backend origin with one backend path.
func BuildBackendURL(baseURL string, path string) (url.URL, error) {
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

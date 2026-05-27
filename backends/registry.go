package backends

import (
	"fmt"
	"os"
	"strings"

	maaslogin "maas-ldap/backends/maas/handlers/login"
)

var definitions = map[string]Definition{
	"maas": {
		Name:            "maas",
		BaseURLEnv:      "MAAS_URL",
		AllowedGroupEnv: "MAAS_LDAP_ALLOWED_GROUP",
		LoginPath:       "/MAAS/accounts/login/",
		NewLoginHandler: maaslogin.NewHandler,
	},
}

// LoadEnabledConfigs loads and validates configuration for each backend listed in BACKENDS.
// Backend names are comma-separated, case-insensitive, and must match entries in definitions.
// It returns an error when BACKENDS is empty, contains duplicates, or names an unknown backend.
func LoadEnabledConfigs() ([]Config, error) {
	backendNames := strings.TrimSpace(os.Getenv("BACKENDS"))
	if backendNames == "" {
		return nil, fmt.Errorf("backend configuration is incomplete. Please set BACKENDS.")
	}

	var configs []Config
	seen := map[string]bool{}
	for _, backendName := range strings.Split(backendNames, ",") {
		backendName = strings.ToLower(strings.TrimSpace(backendName))
		if backendName == "" {
			continue
		}

		if seen[backendName] {
			return nil, fmt.Errorf("backend %q is configured more than once in BACKENDS", backendName)
		}
		seen[backendName] = true

		definition, ok := definitions[backendName]
		if !ok {
			return nil, fmt.Errorf("unknown backend %q in BACKENDS", backendName)
		}

		config, err := LoadConfig(definition)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	if len(configs) == 0 {
		return nil, fmt.Errorf("backend configuration is incomplete. Please set BACKENDS.")
	}

	return configs, nil
}

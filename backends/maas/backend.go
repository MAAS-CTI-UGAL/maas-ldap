package maas

import (
	"fmt"
	"maas-ldap/backends"
	"maas-ldap/users"
	"os"
	"strings"
)

const envMAASURL = "MAAS_URL"
const envMAASAllowedGroup = "MAAS_LDAP_ALLOWED_GROUP"

// Backend contains the MAAS runtime dependencies.
type Backend struct {
	Targets      backends.BackendTargets
	Users        *users.Store
	AllowedGroup string
}

// LoadBackendConfig loads and validates MAAS backend configuration.
func LoadBackendConfig() (backends.BackendTargets, error) {
	return backends.LoadBackendTargets(envMAASURL, Paths)
}

// LoadAllowedGroup loads the LDAP group allowed to access MAAS.
func LoadAllowedGroup() (string, error) {
	allowedGroup := strings.TrimSpace(os.Getenv(envMAASAllowedGroup))
	if allowedGroup == "" {
		return "", fmt.Errorf("MAAS backend configuration is incomplete. Please set %s.", envMAASAllowedGroup)
	}
	return allowedGroup, nil
}

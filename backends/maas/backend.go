package maas

import (
	"maas-ldap/backends"
)

const envMAASURL = "MAAS_URL"
const envMAASAllowedGroup = "MAAS_LDAP_ALLOWED_GROUP"

// LoadBackendConfig loads and validates MAAS backend configuration.
func LoadConfig() (backends.Backend, error) {
	return backends.LoadBackendConfig(envMAASURL, envMAASAllowedGroup, TargetPaths)
}

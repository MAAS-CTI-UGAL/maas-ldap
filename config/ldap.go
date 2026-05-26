package config

import (
	"errors"
	"log"
)

var (
	errMissingLDAPURL       = errors.New("LDAP configuration is incomplete. Please set LDAP_URL.")
	errMissingLDAPUPNSuffix = errors.New("LDAP configuration is incomplete. Please set LDAP_UPN_SUFFIX.")
	errMissingLDAPBaseDN    = errors.New("LDAP configuration is incomplete. Please set LDAP_BASE_DN.")
)

// LDAPConfig contains the environment-driven LDAP connection settings.
type LDAPConfig struct {
	URL        string
	UPN_SUFFIX string
	BASE_DN    string
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

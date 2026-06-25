package config

import "testing"

func TestLoadLDAPConfigReturnsEnvironmentValues(t *testing.T) {
	t.Setenv("LDAP_URL", "ldap://ldap.example.test")
	t.Setenv("LDAP_UPN_SUFFIX", "example.test")
	t.Setenv("LDAP_BASE_DN", "DC=example,DC=test")

	cfg := loadLDAPConfig()

	if cfg.URL != "ldap://ldap.example.test" {
		t.Fatalf("URL = %q, want %q", cfg.URL, "ldap://ldap.example.test")
	}
	if cfg.UPN_SUFFIX != "example.test" {
		t.Fatalf("UPN_SUFFIX = %q, want %q", cfg.UPN_SUFFIX, "example.test")
	}
	if cfg.BASE_DN != "DC=example,DC=test" {
		t.Fatalf("BASE_DN = %q, want %q", cfg.BASE_DN, "DC=example,DC=test")
	}
}

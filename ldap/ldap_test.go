package ldap

import (
	"errors"
	"testing"

	"maas-ldap/config"

	"github.com/go-ldap/ldap/v3"
)

func TestDefaultUserFilterEscapesUsername(t *testing.T) {
	got := DefaultUserFilter(`ali*(ce)\`)
	want := `(&(objectClass=user)(sAMAccountName=ali\2a\28ce\29\5c))`

	if got != want {
		t.Fatalf("DefaultUserFilter() = %q, want %q", got, want)
	}
}

func TestCheckAllowedGroupAcceptsFullDNCaseInsensitive(t *testing.T) {
	entry := ldapEntryWithGroups("CN=MAAS Users,OU=Groups,DC=example,DC=test")

	err := CheckAllowedGroup(entry, "cn=maas users,ou=groups,dc=example,dc=test")
	if err != nil {
		t.Fatalf("CheckAllowedGroup() returned error: %v", err)
	}
}

func TestCheckAllowedGroupAcceptsShortCN(t *testing.T) {
	entry := ldapEntryWithGroups("CN=MAAS Users,OU=Groups,DC=example,DC=test")

	err := CheckAllowedGroup(entry, "maas users")
	if err != nil {
		t.Fatalf("CheckAllowedGroup() returned error: %v", err)
	}
}

func TestCheckAllowedGroupRejectsNonMatchingGroup(t *testing.T) {
	entry := ldapEntryWithGroups("CN=Other Users,OU=Groups,DC=example,DC=test")

	err := CheckAllowedGroup(entry, "maas users")
	if !errors.Is(err, errLDAPGroupCheck) {
		t.Fatalf("CheckAllowedGroup() error = %v, want %v", err, errLDAPGroupCheck)
	}
}

func TestCheckAllowedGroupRejectsPartialDNWhenAllowedGroupContainsEquals(t *testing.T) {
	entry := ldapEntryWithGroups("CN=MAAS Users,OU=Groups,DC=example,DC=test")

	err := CheckAllowedGroup(entry, "CN=MAAS Users")
	if !errors.Is(err, errLDAPGroupCheck) {
		t.Fatalf("CheckAllowedGroup() error = %v, want %v", err, errLDAPGroupCheck)
	}
}

func TestLDAPCredentialValidation(t *testing.T) {
	if err := LdapBind("", "secret", ConfigForTest()); !errors.Is(err, errEmptyCredentials) {
		t.Fatalf("LdapBind() error = %v, want %v", err, errEmptyCredentials)
	}

	if _, err := LdapSearch("alice", "", ConfigForTest(), nil, nil); !errors.Is(err, errEmptyCredentials) {
		t.Fatalf("LdapSearch() error = %v, want %v", err, errEmptyCredentials)
	}
}

func TestLDAPConnectionValidation(t *testing.T) {
	cfg := ConfigForTest()
	cfg.URL = "://bad-url"

	if err := LdapBind("alice", "secret", cfg); !errors.Is(err, errLDAPBindConnect) {
		t.Fatalf("LdapBind() error = %v, want %v", err, errLDAPBindConnect)
	}

	if _, err := LdapSearch("alice", "secret", cfg, nil, nil); !errors.Is(err, errLDAPSearchConnect) {
		t.Fatalf("LdapSearch() error = %v, want %v", err, errLDAPSearchConnect)
	}
}

func ldapEntryWithGroups(groups ...string) *ldap.Entry {
	return &ldap.Entry{
		Attributes: []*ldap.EntryAttribute{
			{Name: "memberOf", Values: groups},
		},
	}
}

func ConfigForTest() config.LDAPConfig {
	return config.LDAPConfig{
		URL:        "ldap://example.test",
		UPN_SUFFIX: "example.test",
		BASE_DN:    "DC=example,DC=test",
	}
}

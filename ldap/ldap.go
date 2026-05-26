package ldap

import (
	"errors"
	"fmt"

	config "maas-ldap/config"

	"github.com/go-ldap/ldap/v3"
)

var (
	errEmptyCredentials  = errors.New("empty username or password")
	errLDAPBindConnect   = errors.New("ldap bind connect failed")
	errLDAPBindAuth      = errors.New("ldap bind authentication failed")
	errLDAPSearchConnect = errors.New("ldap search connect failed")
	errLDAPSearchBind    = errors.New("ldap search bind failed")
	errLDAPSearchQuery   = errors.New("ldap search query failed")
)

// SearchFilterFunc builds the LDAP search filter for a backend-specific user lookup.
type SearchFilterFunc func(username string) string

// LdapBind verifies the supplied username and password against LDAP.
func LdapBind(username, password string, config config.LDAPConfig) error {

	if username == "" || password == "" {
		return errEmptyCredentials
	}

	conn, err := ldap.DialURL(config.URL)
	if err != nil {
		return errLDAPBindConnect
	}
	defer conn.Close()

	// Active Directory commonly accepts username@domain for simple binds.
	bindUser := fmt.Sprintf("%s@%s", username, config.UPN_SUFFIX)

	if err := conn.Bind(bindUser, password); err != nil {
		return errLDAPBindAuth
	}

	return nil
}

// LdapSearch runs a backend-provided LDAP search with the supplied user credentials.
func LdapSearch(username string, password string, config config.LDAPConfig, attributes []string) (*ldap.Entry, error) {

	if username == "" || password == "" {
		return nil, errEmptyCredentials
	}

	conn, err := ldap.DialURL(config.URL)
	if err != nil {
		return nil, errLDAPSearchConnect
	}
	defer conn.Close()

	// Search with the same user credentials that were submitted for login.
	bindUser := fmt.Sprintf("%s@%s", username, config.UPN_SUFFIX)
	if err := conn.Bind(bindUser, password); err != nil {
		return nil, errLDAPSearchBind
	}

	req := ldap.NewSearchRequest(
		config.BASE_DN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		1,
		0,
		false,
		LDAPFilter(username),
		attributes,
		nil,
	)

	res, err := conn.Search(req)
	if err != nil {
		return nil, errLDAPSearchQuery
	}

	if len(res.Entries) == 0 {
		return nil, fmt.Errorf("%w: no entries found for user %s", errLDAPSearchQuery, username)
	}

	if len(res.Entries) > 1 {
		return nil, fmt.Errorf("%w: multiple entries found for user %s", errLDAPSearchQuery, username)
	}

	return res.Entries[0], nil
}

func LDAPFilter(username string) string {
	return fmt.Sprintf(
		"(&(objectClass=user)(sAMAccountName=%s))",
		ldap.EscapeFilter(username),
	)
}

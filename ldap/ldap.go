package ldap

import (
	"errors"
	"fmt"
	"strings"

	maasconfig "maas-ldap/config"

	"github.com/go-ldap/ldap/v3"
)

var (
	errMissingLDAPURL       = errors.New("LDAP_URL is missing")
	errMissingLDAPUPNSuffix = errors.New("LDAP_UPN_SUFFIX is missing")
	errMissingLDAPBaseDN    = errors.New("LDAP_BASE_DN is missing")
	errEmptyCredentials     = errors.New("empty username or password")
)

// LdapBind verifies the supplied username and password against LDAP.
func LdapBind(username, password string, config maasconfig.LDAPConfig) error {

	if config.URL == "" {
		return errMissingLDAPURL
	}

	if config.UPN_SUFFIX == "" {
		return errMissingLDAPUPNSuffix
	}

	if username == "" || password == "" {
		return errEmptyCredentials
	}

	conn, err := ldap.DialURL(config.URL)
	if err != nil {
		return fmt.Errorf("ldap connect: %w", err)
	}
	defer conn.Close()

	// Active Directory commonly accepts username@domain for simple binds.
	bindUser := fmt.Sprintf("%s@%s", username, config.UPN_SUFFIX)

	if err := conn.Bind(bindUser, password); err != nil {
		return fmt.Errorf("ldap bind: %w", err)
	}

	return nil
}

// LdapSearch finds the LDAP user and reports whether it belongs to the allowed group.
func LdapSearch(username string, config maasconfig.LDAPConfig) (bool, error) {

	if config.URL == "" {
		return false, errMissingLDAPURL
	}

	if config.UPN_SUFFIX == "" {
		return false, errMissingLDAPUPNSuffix
	}

	if config.BASE_DN == "" {
		return false, errMissingLDAPBaseDN
	}

	conn, err := ldap.DialURL(config.URL)
	if err != nil {
		return false, fmt.Errorf("ldap connect: %w", err)
	}
	defer conn.Close()

	// Escape the username before it is embedded in the LDAP filter.
	filter := fmt.Sprintf(
		"(&(objectClass=user)(sAMAccountName=%s))",
		ldap.EscapeFilter(username),
	)

	req := ldap.NewSearchRequest(
		config.BASE_DN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		1,
		0,
		false,
		filter,
		[]string{"memberOf"},
		nil,
	)

	res, err := conn.Search(req)
	if err != nil {
		return false, fmt.Errorf("ldap search: %w", err)
	}

	if len(res.Entries) != 1 {
		return false, fmt.Errorf("expected 1 user, got %d", len(res.Entries))
	}

	// Membership values are full DNs. LDAP_ALLOWED_GROUP can be either a full DN
	// or a short CN value.
	for _, group := range res.Entries[0].GetAttributeValues("memberOf") {
		if isAllowedGroup(group, config.ALLOWED_GROUP) {
			return true, nil
		}
	}

	return false, nil
}

func isAllowedGroup(memberOf string, allowedGroup string) bool {
	if strings.EqualFold(memberOf, allowedGroup) {
		return true
	}

	if strings.Contains(allowedGroup, "=") {
		return false
	}

	return strings.HasPrefix(
		strings.ToLower(memberOf),
		"cn="+strings.ToLower(allowedGroup)+",",
	)
}

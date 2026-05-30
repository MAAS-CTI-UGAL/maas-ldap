package maas

import (
	"errors"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

var (
	errLDAPMissingMAASPassword = errors.New("ldap entry is missing maas password")
	errLDAPGroupCheck          = errors.New("user is not in allowed group")
)

func checkAllowedGroup(entry *ldap.Entry, allowedGroup string) bool {
	// Membership values are full DNs. The allowed group can be either a full DN
	// or a short CN value.
	for _, group := range entry.GetAttributeValues("memberOf") {
		if isAllowedGroup(group, allowedGroup) {
			return true
		}
	}

	return false
}

func maasPassword(entry *ldap.Entry) (string, error) {
	values := entry.GetAttributeValues("primaryTelexNumber")
	if len(values) != 1 || strings.TrimSpace(values[0]) == "" {
		return "", errLDAPMissingMAASPassword
	}
	return values[0], nil
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

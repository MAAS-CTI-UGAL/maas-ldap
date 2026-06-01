package maas

import (
	"errors"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

var errLDAPMissingMAASPassword = errors.New("ldap entry is missing maas password")

func maasPassword(entry *ldap.Entry) (string, error) {
	values := entry.GetAttributeValues("primaryTelexNumber")
	if len(values) != 1 || strings.TrimSpace(values[0]) == "" {
		return "", errLDAPMissingMAASPassword
	}
	return values[0], nil
}

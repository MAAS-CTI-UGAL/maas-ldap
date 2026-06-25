package maas

import (
	"errors"
	"testing"

	"github.com/go-ldap/ldap/v3"
)

func TestMAASPasswordReturnsSingleNonBlankValue(t *testing.T) {
	entry := &ldap.Entry{
		Attributes: []*ldap.EntryAttribute{
			{Name: "primaryTelexNumber", Values: []string{"maas-secret"}},
		},
	}

	got, err := maasPassword(entry)
	if err != nil {
		t.Fatalf("maasPassword() returned error: %v", err)
	}
	if got != "maas-secret" {
		t.Fatalf("maasPassword() = %q, want %q", got, "maas-secret")
	}
}

func TestMAASPasswordValidation(t *testing.T) {
	tests := []struct {
		name  string
		entry *ldap.Entry
	}{
		{
			name:  "missing",
			entry: &ldap.Entry{},
		},
		{
			name: "blank",
			entry: &ldap.Entry{
				Attributes: []*ldap.EntryAttribute{
					{Name: "primaryTelexNumber", Values: []string{"  "}},
				},
			},
		},
		{
			name: "multiple",
			entry: &ldap.Entry{
				Attributes: []*ldap.EntryAttribute{
					{Name: "primaryTelexNumber", Values: []string{"one", "two"}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := maasPassword(tt.entry)
			if !errors.Is(err, errLDAPMissingMAASPassword) {
				t.Fatalf("maasPassword() error = %v, want %v", err, errLDAPMissingMAASPassword)
			}
		})
	}
}

package ldap

import "testing"

func TestIsAllowedGroup(t *testing.T) {
	tests := []struct {
		name         string
		memberOf     string
		allowedGroup string
		want         bool
	}{
		{
			name:         "full DN exact match",
			memberOf:     "CN=MaaS_Allowed,OU=Groups,DC=example,DC=internal",
			allowedGroup: "CN=MaaS_Allowed,OU=Groups,DC=example,DC=internal",
			want:         true,
		},
		{
			name:         "full DN case insensitive match",
			memberOf:     "CN=MaaS_Allowed,OU=Groups,DC=example,DC=internal",
			allowedGroup: "cn=maas_allowed,ou=groups,dc=example,dc=internal",
			want:         true,
		},
		{
			name:         "short CN match",
			memberOf:     "CN=MaaS_Allowed,OU=Groups,DC=example,DC=internal",
			allowedGroup: "MaaS_Allowed",
			want:         true,
		},
		{
			name:         "short CN case insensitive match",
			memberOf:     "CN=MaaS_Allowed,OU=Groups,DC=example,DC=internal",
			allowedGroup: "maas_allowed",
			want:         true,
		},
		{
			name:         "short CN does not match longer group name",
			memberOf:     "CN=MaaS_Allowed_ReadOnly,OU=Groups,DC=example,DC=internal",
			allowedGroup: "MaaS_Allowed",
			want:         false,
		},
		{
			name:         "short CN does not match unrelated group",
			memberOf:     "CN=Other_Group,OU=Groups,DC=example,DC=internal",
			allowedGroup: "MaaS_Allowed",
			want:         false,
		},
		{
			name:         "full DN does not allow prefix-only match",
			memberOf:     "CN=MaaS_Allowed,OU=Groups,DC=example,DC=internal",
			allowedGroup: "CN=MaaS_Allowed,OU=Other,DC=example,DC=internal",
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isAllowedGroup(tt.memberOf, tt.allowedGroup)
			if got != tt.want {
				t.Fatalf("isAllowedGroup(%q, %q) = %v, want %v", tt.memberOf, tt.allowedGroup, got, tt.want)
			}
		})
	}
}

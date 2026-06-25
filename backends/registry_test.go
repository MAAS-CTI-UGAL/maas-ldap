package backends

import (
	"strings"
	"testing"
)

func TestLoadEnabledConfigsLoadsConfiguredBackends(t *testing.T) {
	t.Setenv("BACKENDS", " MAAS , maas-manager ")
	t.Setenv("MAAS_URL", "https://maas.example.test/root")
	t.Setenv("MAAS_LDAP_ALLOWED_GROUP", "maas-users")
	t.Setenv("MAAS_MANAGER_URL", "https://manager.example.test")
	t.Setenv("MAAS_MANAGER_LDAP_ALLOWED_GROUP", "manager-users")

	configs, err := LoadEnabledConfigs()
	if err != nil {
		t.Fatalf("LoadEnabledConfigs() returned error: %v", err)
	}
	if len(configs) != 2 {
		t.Fatalf("len(configs) = %d, want 2", len(configs))
	}
	if configs[0].Name != "maas" || configs[0].Target.String() != "https://maas.example.test/root/MAAS/accounts/login/" {
		t.Fatalf("first config = %#v", configs[0])
	}
	if configs[1].Name != "maas-manager" || configs[1].Target.String() != "https://manager.example.test/manager/api/login" {
		t.Fatalf("second config = %#v", configs[1])
	}
}

func TestLoadEnabledConfigsValidation(t *testing.T) {
	tests := []struct {
		name     string
		backends string
		wantErr  string
	}{
		{
			name:     "missing",
			backends: "",
			wantErr:  "backend configuration is incomplete. Please set BACKENDS.",
		},
		{
			name:     "blank entries only",
			backends: ", ,",
			wantErr:  "backend configuration is incomplete. Please set BACKENDS.",
		},
		{
			name:     "duplicate",
			backends: "maas, MAAS",
			wantErr:  `backend "maas" is configured more than once in BACKENDS`,
		},
		{
			name:     "unknown",
			backends: "unknown",
			wantErr:  `unknown backend "unknown" in BACKENDS`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("BACKENDS", tt.backends)
			t.Setenv("MAAS_URL", "https://maas.example.test")
			t.Setenv("MAAS_LDAP_ALLOWED_GROUP", "maas-users")

			_, err := LoadEnabledConfigs()
			if err == nil {
				t.Fatal("LoadEnabledConfigs() returned nil error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("LoadEnabledConfigs() error = %q, want containing %q", err, tt.wantErr)
			}
		})
	}
}

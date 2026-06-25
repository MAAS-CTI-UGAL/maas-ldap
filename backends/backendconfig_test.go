package backends

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"maas-ldap/config"
)

func stubLoginHandler(config.AppConfig, url.URL, string) http.HandlerFunc {
	return func(http.ResponseWriter, *http.Request) {}
}

func TestLoadConfigBuildsTargetAndTrimsAllowedGroup(t *testing.T) {
	t.Setenv("TEST_BACKEND_URL", "https://example.test/base/")
	t.Setenv("TEST_ALLOWED_GROUP", " allowed ")

	cfg, err := LoadBackendConfig(BackendDefinition{
		Name:            "test",
		BaseURLEnv:      "TEST_BACKEND_URL",
		AllowedGroupEnv: "TEST_ALLOWED_GROUP",
		LoginPath:       "/login/",
		NewLoginHandler: stubLoginHandler,
	})
	if err != nil {
		t.Fatalf("LoadBackendConfig() returned error: %v", err)
	}

	if cfg.BaseURL != "https://example.test/base/" {
		t.Fatalf("BaseURL = %q, want %q", cfg.BaseURL, "https://example.test/base/")
	}
	if cfg.Target.String() != "https://example.test/base/login/" {
		t.Fatalf("Target = %q, want %q", cfg.Target.String(), "https://example.test/base/login/")
	}
	if cfg.AllowedGroup != "allowed" {
		t.Fatalf("AllowedGroup = %q, want %q", cfg.AllowedGroup, "allowed")
	}
}

func TestLoadConfigRequiresLoginHandler(t *testing.T) {
	_, err := LoadBackendConfig(BackendDefinition{Name: "broken"})
	if err == nil {
		t.Fatal("LoadBackendConfig() returned nil error")
	}
	if !strings.Contains(err.Error(), `backend "broken" login handler is not configured`) {
		t.Fatalf("LoadBackendConfig() error = %q", err)
	}
}

func TestLoadBackendTargetValidation(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr string
	}{
		{
			name:    "missing",
			value:   "",
			wantErr: "backend configuration is incomplete. Please set TEST_BACKEND_URL.",
		},
		{
			name:    "missing scheme",
			value:   "example.test",
			wantErr: "TEST_BACKEND_URL must include scheme and host",
		},
		{
			name:    "missing host",
			value:   "https:///login",
			wantErr: "TEST_BACKEND_URL must include scheme and host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TEST_BACKEND_URL", tt.value)

			_, _, err := loadBackendTarget("TEST_BACKEND_URL", "/login")
			if err == nil {
				t.Fatal("loadBackendTarget() returned nil error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("loadBackendTarget() error = %q, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestLoadAllowedGroupValidation(t *testing.T) {
	t.Setenv("TEST_ALLOWED_GROUP", "")

	_, err := LoadAllowedGroup("TEST_ALLOWED_GROUP")
	if err == nil {
		t.Fatal("LoadAllowedGroup() returned nil error")
	}
	if !strings.Contains(err.Error(), "backend configuration is incomplete. Please set TEST_ALLOWED_GROUP.") {
		t.Fatalf("LoadAllowedGroup() error = %q", err)
	}
}

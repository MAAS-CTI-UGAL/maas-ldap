package maas_manager

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeLoginRequestAcceptsJSONAndTrimsUsername(t *testing.T) {
	req := httptest.NewRequest("POST", "/manager/api/login", strings.NewReader(`{"username":" alice ","password":"secret"}`))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	got, err := decodeLoginRequest(req)
	if err != nil {
		t.Fatalf("decodeLoginRequest() returned error: %v", err)
	}
	if got.Username != "alice" {
		t.Fatalf("Username = %q, want %q", got.Username, "alice")
	}
	if got.Password != "secret" {
		t.Fatalf("Password = %q, want %q", got.Password, "secret")
	}
}

func TestDecodeLoginRequestValidation(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
		wantErr     error
	}{
		{
			name:        "unexpected content type",
			contentType: "application/x-www-form-urlencoded",
			body:        `{"username":"alice","password":"secret"}`,
			wantErr:     errUnexpectedContentType,
		},
		{
			name:        "missing username",
			contentType: "application/json",
			body:        `{"username":" ","password":"secret"}`,
			wantErr:     errEmptyUsername,
		},
		{
			name:        "missing password",
			contentType: "application/json",
			body:        `{"username":"alice","password":" "}`,
			wantErr:     errEmptyPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/manager/api/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)

			_, err := decodeLoginRequest(req)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("decodeLoginRequest() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodeLoginRequestRejectsInvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/manager/api/login", strings.NewReader(`{"username":`))
	req.Header.Set("Content-Type", "application/json")

	if _, err := decodeLoginRequest(req); err == nil {
		t.Fatal("decodeLoginRequest() returned nil error")
	}
}

func TestIsJSONContentType(t *testing.T) {
	tests := []struct {
		contentType string
		want        bool
	}{
		{contentType: "application/json", want: true},
		{contentType: "Application/JSON; charset=utf-8", want: true},
		{contentType: "application/x-www-form-urlencoded", want: false},
		{contentType: "%%%invalid", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.contentType, func(t *testing.T) {
			if got := isJSONContentType(tt.contentType); got != tt.want {
				t.Fatalf("isJSONContentType(%q) = %v, want %v", tt.contentType, got, tt.want)
			}
		})
	}
}

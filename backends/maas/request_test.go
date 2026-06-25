package maas

import (
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeLoginRequestAcceptsFormAndPreservesBody(t *testing.T) {
	body := "username=alice&password=secret&next=%2FMAAS"
	req := httptest.NewRequest("POST", "/MAAS/accounts/login/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	form, err := decodeLoginRequest(req)
	if err != nil {
		t.Fatalf("decodeLoginRequest() returned error: %v", err)
	}
	if form.Get("username") != "alice" {
		t.Fatalf("username = %q, want %q", form.Get("username"), "alice")
	}
	if form.Get("password") != "secret" {
		t.Fatalf("password = %q, want %q", form.Get("password"), "secret")
	}

	replayed, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("ReadAll(replayed body) returned error: %v", err)
	}
	if string(replayed) != body {
		t.Fatalf("replayed body = %q, want %q", string(replayed), body)
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
			contentType: "application/json",
			body:        "username=alice&password=secret",
			wantErr:     errUnexpectedContentType,
		},
		{
			name:        "missing username",
			contentType: "application/x-www-form-urlencoded",
			body:        "password=secret",
			wantErr:     errEmptyUsername,
		},
		{
			name:        "missing password",
			contentType: "application/x-www-form-urlencoded",
			body:        "username=alice",
			wantErr:     errEmptyPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/MAAS/accounts/login/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)

			_, err := decodeLoginRequest(req)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("decodeLoginRequest() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsFormContentType(t *testing.T) {
	tests := []struct {
		contentType string
		want        bool
	}{
		{contentType: "application/x-www-form-urlencoded", want: true},
		{contentType: "Application/X-WWW-Form-Urlencoded; charset=utf-8", want: true},
		{contentType: "application/json", want: false},
		{contentType: "%%%invalid", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.contentType, func(t *testing.T) {
			if got := isFormContentType(tt.contentType); got != tt.want {
				t.Fatalf("isFormContentType(%q) = %v, want %v", tt.contentType, got, tt.want)
			}
		})
	}
}

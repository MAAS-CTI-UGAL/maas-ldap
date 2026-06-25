package maas

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"maas-ldap/config"
)

func TestNewHandlerRejectsNonPostRequests(t *testing.T) {
	var got operationError
	restoreWriteError(t, &got)

	handler := NewHandler(config.AppConfig{}, url.URL{}, "allowed")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/MAAS/accounts/login/", nil)

	handler(rr, req)

	if rr.Header().Get("Allow") != http.MethodPost {
		t.Fatalf("Allow = %q, want %q", rr.Header().Get("Allow"), http.MethodPost)
	}
	assertOperationError(t, got, "Invalid HTTP method", "This page only accepts login submissions.", nil, http.StatusMethodNotAllowed)
}

func TestNewHandlerRejectsInvalidLoginRequest(t *testing.T) {
	var got operationError
	restoreWriteError(t, &got)

	handler := NewHandler(config.AppConfig{}, url.URL{}, "allowed")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/MAAS/accounts/login/", strings.NewReader("username=alice"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(rr, req)

	assertOperationError(t, got, "Invalid login request", "Please submit the login form again.", errEmptyPassword, http.StatusBadRequest)
}

func TestNewHandlerReportsLDAPSearchFailure(t *testing.T) {
	var got operationError
	restoreWriteError(t, &got)

	handler := NewHandler(config.AppConfig{
		LDAP: config.LDAPConfig{
			URL:        "://bad-url",
			UPN_SUFFIX: "example.test",
			BASE_DN:    "DC=example,DC=test",
		},
	}, url.URL{}, "allowed")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/MAAS/accounts/login/", strings.NewReader("username=alice&password=secret"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(rr, req)

	if got.operation != "/MAAS/accounts/login/" {
		t.Fatalf("operation = %q, want %q", got.operation, "/MAAS/accounts/login/")
	}
	if got.logMessage != "LDAP search failed" {
		t.Fatalf("logMessage = %q, want %q", got.logMessage, "LDAP search failed")
	}
	if got.userMessage != "We could not verify your MAAS access. Please try again or contact an administrator." {
		t.Fatalf("userMessage = %q", got.userMessage)
	}
	if got.err == nil {
		t.Fatal("err is nil")
	}
	if got.statusCode != http.StatusUnauthorized {
		t.Fatalf("statusCode = %d, want %d", got.statusCode, http.StatusUnauthorized)
	}
}

type operationError struct {
	operation   string
	logMessage  string
	userMessage string
	err         error
	statusCode  int
}

func restoreWriteError(t *testing.T, got *operationError) {
	t.Helper()

	original := WriteError
	WriteError = func(_ http.ResponseWriter, operation string, logMessage string, userMessage string, err error, statusCode int) {
		*got = operationError{
			operation:   operation,
			logMessage:  logMessage,
			userMessage: userMessage,
			err:         err,
			statusCode:  statusCode,
		}
	}
	t.Cleanup(func() { WriteError = original })
}

func assertOperationError(t *testing.T, got operationError, logMessage string, userMessage string, wantErr error, statusCode int) {
	t.Helper()

	if got.operation != "/MAAS/accounts/login/" {
		t.Fatalf("operation = %q, want %q", got.operation, "/MAAS/accounts/login/")
	}
	if got.logMessage != logMessage {
		t.Fatalf("logMessage = %q, want %q", got.logMessage, logMessage)
	}
	if got.userMessage != userMessage {
		t.Fatalf("userMessage = %q, want %q", got.userMessage, userMessage)
	}
	if !errors.Is(got.err, wantErr) {
		t.Fatalf("err = %v, want %v", got.err, wantErr)
	}
	if got.statusCode != statusCode {
		t.Fatalf("statusCode = %d, want %d", got.statusCode, statusCode)
	}
}

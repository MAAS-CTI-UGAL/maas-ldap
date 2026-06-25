package errorwriter

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestErrorWriterWritesPublicMessageAndLogsInternalError(t *testing.T) {
	var logs bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(originalOutput) })

	writer := New("maas")
	rr := httptest.NewRecorder()

	writer(rr, "/login", "LDAP search failed", "public message", errors.New("bind failed"), http.StatusUnauthorized)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
	if rr.Body.String() != "public message\n" {
		t.Fatalf("body = %q, want %q", rr.Body.String(), "public message\n")
	}
	if !strings.Contains(logs.String(), "maas backend /login failed: LDAP search failed: bind failed") {
		t.Fatalf("logs = %q", logs.String())
	}
}

func TestErrorWriterLogsWithoutInternalError(t *testing.T) {
	var logs bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(originalOutput) })

	writer := New("maas-manager")
	rr := httptest.NewRecorder()

	writer(rr, "/login", "Invalid login request", "try again", nil, http.StatusBadRequest)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
	if !strings.Contains(logs.String(), "maas-manager backend /login failed: Invalid login request") {
		t.Fatalf("logs = %q", logs.String())
	}
	if strings.Contains(logs.String(), "bind failed") {
		t.Fatalf("logs contain unexpected internal error: %q", logs.String())
	}
}

func TestMain(m *testing.M) {
	log.SetFlags(0)
	os.Exit(m.Run())
}

package global_handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddRoutesRegistersHealth(t *testing.T) {
	mux := http.NewServeMux()
	AddRoutes(mux)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, Health, nil)

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if rr.Body.String() != "ok" {
		t.Fatalf("body = %q, want %q", rr.Body.String(), "ok")
	}
}

package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"maas-ldap/backends"
	"maas-ldap/config"
)

func TestAddRoutesRegistersGlobalAndBackendRoutes(t *testing.T) {
	mux := http.NewServeMux()
	AddRoutes(mux, config.AppConfig{}, []backends.BackendConfig{
		{
			BackendDefinition: backends.BackendDefinition{
				LoginPath: "/backend/login",
				NewLoginHandler: func(config.AppConfig, url.URL, string) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusAccepted)
						_, _ = w.Write([]byte("backend"))
					}
				},
			},
		},
	})

	assertRoute(t, mux, http.MethodGet, "/health", http.StatusOK, "ok")
	assertRoute(t, mux, http.MethodPost, "/backend/login", http.StatusAccepted, "backend")
}

func assertRoute(t *testing.T, mux *http.ServeMux, method string, path string, status int, body string) {
	t.Helper()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)

	mux.ServeHTTP(rr, req)

	if rr.Code != status {
		t.Fatalf("%s status = %d, want %d", path, rr.Code, status)
	}
	if rr.Body.String() != body {
		t.Fatalf("%s body = %q, want %q", path, rr.Body.String(), body)
	}
}

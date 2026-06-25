package backends

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"maas-ldap/config"
)

func TestAddRoutesRegistersBackendHandlers(t *testing.T) {
	mux := http.NewServeMux()
	AddRoutes(mux, config.AppConfig{}, []BackendConfig{
		{
			BackendDefinition: BackendDefinition{
				LoginPath: "/backend/login",
				NewLoginHandler: func(config.AppConfig, url.URL, string) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusCreated)
						_, _ = w.Write([]byte("backend"))
					}
				},
			},
		},
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/backend/login", nil)

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusCreated)
	}
	if rr.Body.String() != "backend" {
		t.Fatalf("body = %q, want %q", rr.Body.String(), "backend")
	}
}

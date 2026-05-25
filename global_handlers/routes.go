package global_handlers

import (
	"maas-ldap/global_handlers/health"
	"net/http"
)

// AddRoutes registers health check routes.
func AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc(Health, health.NewHandler())
}

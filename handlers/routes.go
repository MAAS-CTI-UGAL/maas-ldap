package handlers

import (
	"net/http"

	"maas-ldap/config"
	"maas-ldap/handlers/health"
	"maas-ldap/handlers/login"
)

// AddRoutes registers all HTTP routes exposed by the proxy.
func AddRoutes(mux *http.ServeMux, appConfig config.AppConfig) {
	health.AddRoutes(mux)
	login.AddRoutes(mux, appConfig)
}

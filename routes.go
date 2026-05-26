package main

import (
	"net/http"

	"maas-ldap/backends"
	"maas-ldap/backends/maas"
	"maas-ldap/config"
	handlers "maas-ldap/global_handlers"
)

// AddRoutes registers all HTTP routes exposed by the proxy.
func AddRoutes(mux *http.ServeMux, appConfig config.AppConfig, maasBackend backends.Backend) {
	handlers.AddRoutes(mux)
	maas.AddRoutes(mux, appConfig, maasBackend)
}

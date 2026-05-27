package main

import (
	"net/http"

	"maas-ldap/backends"
	"maas-ldap/config"
	global_handlers "maas-ldap/global_handlers"
)

// AddRoutes registers all HTTP routes exposed by the proxy.
func AddRoutes(mux *http.ServeMux, appConfig config.AppConfig, backendConfigs []backends.Config) {
	global_handlers.AddRoutes(mux)
	backends.AddRoutes(mux, appConfig, backendConfigs)
}

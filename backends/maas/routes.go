package maas

import (
	"net/http"

	"maas-ldap/backends"
	"maas-ldap/backends/maas/handlers/login"
	"maas-ldap/config"
)

// AddRoutes registers login routes.
func AddRoutes(mux *http.ServeMux, appConfig config.AppConfig, maas backends.Backend) {
	mux.HandleFunc(LoginPath, login.NewHandler(appConfig, maas.Targets[LoginTarget], maas.AllowedGroup))
}

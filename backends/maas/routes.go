package maas

import (
	"net/http"

	maaslogin "maas-ldap/backends/maas/handlers/login"
	"maas-ldap/config"
)

// AddRoutes registers login routes.
func AddRoutes(mux *http.ServeMux, appConfig config.AppConfig, backend Backend) {
	mux.HandleFunc(LoginPath, maaslogin.NewHandler(appConfig, backend.Users, backend.Targets.Targets[LoginTarget], backend.AllowedGroup))
}

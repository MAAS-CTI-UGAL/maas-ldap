package login

import (
	"net/http"

	"maas-ldap/backends/maas"
	"maas-ldap/config"
)

// AddRoutes registers login routes.
func AddRoutes(mux *http.ServeMux, appConfig config.AppConfig) {
	mux.HandleFunc(maas.LoginPath, NewHandler(appConfig))
}

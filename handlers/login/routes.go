package login

import (
	"net/http"

	"maas-ldap/config"
)

// AddRoutes registers login routes.
func AddRoutes(mux *http.ServeMux, appConfig config.AppConfig) {
	mux.HandleFunc(config.LoginMAASPath, NewHandler(appConfig))
}

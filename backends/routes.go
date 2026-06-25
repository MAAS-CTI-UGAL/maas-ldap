package backends

import (
	"net/http"

	"maas-ldap/config"
)

// AddRoutes registers all enabled backend login routes.
func AddRoutes(mux *http.ServeMux, appConfig config.AppConfig, backendConfigs []BackendConfig) {
	for _, backendConfig := range backendConfigs {
		mux.HandleFunc(
			backendConfig.LoginPath,
			backendConfig.NewLoginHandler(appConfig, backendConfig.Target, backendConfig.AllowedGroup),
		)
	}
}

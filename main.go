package main

import (
	"log"
	"maas-ldap/backends"
	"maas-ldap/config"
	"maas-ldap/middlewares"
	"net/http"
)

func main() {
	// Load environment configuration before wiring handlers so startup fails fast on bad env.
	appConfig := config.Bootstrap()

	backendConfigs, err := backends.LoadEnabledConfigs()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	AddRoutes(mux, appConfig, backendConfigs)

	log.Printf("Server running on %s", appConfig.ListenAddress)
	err = http.ListenAndServe(appConfig.ListenAddress, middlewares.Logging(mux))
	if err != nil {
		log.Fatal(err)
	}
}

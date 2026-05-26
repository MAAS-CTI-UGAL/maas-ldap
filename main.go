package main

import (
	"log"
	"maas-ldap/backends/maas"
	"maas-ldap/config"
	"maas-ldap/global_handlers"
	"maas-ldap/middlewares"
	"net/http"
)

func main() {
	// Load environment configuration before wiring handlers so startup fails fast on bad env.
	appConfig := config.Bootstrap()

	maasBackend, err := maas.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	maas.AddRoutes(mux, appConfig, maasBackend)
	global_handlers.AddRoutes(mux)

	log.Printf("Server running on %s", appConfig.ListenAddress)
	err = http.ListenAndServe(appConfig.ListenAddress, middlewares.Logging(mux))
	if err != nil {
		log.Fatal(err)
	}
}

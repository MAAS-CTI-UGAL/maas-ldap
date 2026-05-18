package main

import (
	"log"
	"maas-ldap/config"
	"maas-ldap/handlers"
	"maas-ldap/logging"
	"net/http"
)

func main() {
	// Load configuration before wiring handlers so startup fails fast on bad env.
	appConfig := config.Bootstrap()
	logFile, err := logging.Configure(appConfig.App.LogFile)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	mux := http.NewServeMux()
	handlers.AddRoutes(mux, appConfig)

	log.Printf("Server running on %s", appConfig.App.ListenAddress)
	err = http.ListenAndServe(appConfig.App.ListenAddress, mux)
	if err != nil {
		log.Fatal(err)
	}
}

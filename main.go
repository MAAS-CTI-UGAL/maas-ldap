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
	logFile, err := logging.Configure(appConfig.Settings.LogFilePath)
	if err != nil {
		log.Fatal(err)
	}
	if logFile != nil {
		defer logFile.Close()
	}

	mux := http.NewServeMux()
	handlers.AddRoutes(mux, appConfig)

	log.Printf("Server running on %s", appConfig.Settings.ListenAddress)
	err = http.ListenAndServe(appConfig.Settings.ListenAddress, mux)
	if err != nil {
		log.Fatal(err)
	}
}

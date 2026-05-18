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

	loginHandler := handlers.NewLoginHandler(appConfig)

	// The proxy only exposes the target login route.
	http.HandleFunc(appConfig.App.LoginPath, loginHandler)

	log.Printf("Server running on %s", appConfig.App.ListenAddress)
	err = http.ListenAndServe(appConfig.App.ListenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}

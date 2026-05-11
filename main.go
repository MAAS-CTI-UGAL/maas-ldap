package main

import (
	"log"
	"maas-ldap/config"
	"maas-ldap/handlers"
	"maas-ldap/logging"
	"net/http"
)

func main() {
	appConfig := config.Bootstrap()
	logFile, err := logging.Configure(appConfig.App.LogFile)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	loginHandler := handlers.NewLoginHandler(appConfig)

	http.HandleFunc(appConfig.App.LoginPath, loginHandler)

	log.Printf("Server running on %s", appConfig.App.ListenAddress)
	err = http.ListenAndServe(appConfig.App.ListenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}

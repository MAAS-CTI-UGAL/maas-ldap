package main

import (
	"log"
	"maas-ldap/backends/maas"
	"maas-ldap/config"
	"maas-ldap/db"
	"maas-ldap/logging"
	"maas-ldap/users"
	"net/http"
)

func main() {
	// Load environment configuration before wiring handlers so startup fails fast on bad env.
	appConfig := config.Bootstrap()
	logFile, err := logging.Configure(appConfig.Settings.LogFilePath)
	if err != nil {
		log.Fatal(err)
	}
	if logFile != nil {
		defer logFile.Close()
	}

	database, err := db.Open(appConfig.Settings.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	maasConfig, err := maas.LoadBackendConfig()
	if err != nil {
		log.Fatal(err)
	}

	maasAllowedGroup, err := maas.LoadAllowedGroup()
	if err != nil {
		log.Fatal(err)
	}

	maasMappings, err := maas.LoadUserMappings(database)
	if err != nil {
		log.Fatal(err)
	}
	maasBackend := maas.Backend{
		Targets:      maasConfig,
		Users:        users.NewStore(maasMappings),
		AllowedGroup: maasAllowedGroup,
	}

	mux := http.NewServeMux()
	maas.AddRoutes(mux, appConfig, maasBackend)

	log.Printf("Server running on %s", appConfig.Settings.ListenAddress)
	err = http.ListenAndServe(appConfig.Settings.ListenAddress, mux)
	if err != nil {
		log.Fatal(err)
	}
}

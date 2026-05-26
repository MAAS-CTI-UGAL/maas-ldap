package main

import (
	"log"
	"maas-ldap/backends/maas"
	"maas-ldap/config"
	"maas-ldap/db"
	"maas-ldap/global_handlers"
	"maas-ldap/middlewares"
	"maas-ldap/users"
	"net/http"
)

func main() {
	// Load environment configuration before wiring handlers so startup fails fast on bad env.
	appConfig := config.Bootstrap()
	if appConfig.Settings.Log.File != nil {
		defer appConfig.Settings.Log.File.Close()
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
	global_handlers.AddRoutes(mux)

	log.Printf("Server running on %s", appConfig.Settings.ListenAddress)
	err = http.ListenAndServe(appConfig.Settings.ListenAddress, middlewares.Logging(mux))
	if err != nil {
		log.Fatal(err)
	}
}

package config

import "os"

const (
	defaultListenAddress = ":42069"
	defaultLogFile       = "/var/log/maas-ldap.log"
	envLogFile           = "MAAS_LDAP_LOG_FILE"
)

// AppSettings contains application-wide routes and file paths.
type AppSettings struct {
	ListenAddress string
	LoginPath     string
	LogFile       string
	UsersFile     string
}

func loadAppSettings() AppSettings {
	return AppSettings{
		ListenAddress: defaultListenAddress,
		LoginPath:     defaultLoginPath,
		LogFile:       envOrDefault(envLogFile, defaultLogFile),
		UsersFile:     defaultUsersFile,
	}
}

// envOrDefault lets operators override selected defaults without changing code.
func envOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

package config

const (
	// The local mapping file supplies MAAS passwords for LDAP users.
	defaultUsersFile = "users.json"
)

const (
	// LoginMAAS identifies the MAAS login backend endpoint.
	LoginMAAS = "loginMAAS"
	// LoginMAASPath is the local and backend path MAAS posts credentials to.
	LoginMAASPath = "/MAAS/accounts/login/"
)

// BackendEndpointPaths defines each backend endpoint path by lookup key.
var BackendEndpointPaths = map[string]string{
	LoginMAAS: LoginMAASPath,
}

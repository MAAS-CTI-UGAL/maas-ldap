package config

const (
	// MAAS posts credentials to this path during interactive login.
	defaultLoginPath = "/MAAS/accounts/login/"
	// The local mapping file supplies MAAS passwords for LDAP users.
	defaultUsersFile = "users.json"
)

// BackendEndpointPaths defines each backend endpoint path by lookup key.
var BackendEndpointPaths = map[string]string{
	"login": defaultLoginPath,
}

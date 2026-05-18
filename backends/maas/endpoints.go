package maas

const (
	// LoginEndpoint identifies the MAAS login backend endpoint.
	LoginEndpoint = "login"
	// LoginPath is the path MAAS posts credentials to during interactive login.
	LoginPath = "/MAAS/accounts/login/"
)

// EndpointPaths defines MAAS backend endpoint paths by lookup key.
var EndpointPaths = map[string]string{
	LoginEndpoint: LoginPath,
}

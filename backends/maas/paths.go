package maas

const (
	LoginTarget  = "login"
	LogoutTarget = "logout"

	LoginPath  = "/MAAS/accounts/login/"
	LogoutPath = "/MAAS/accounts/logout/"
)

var Paths = map[string]string{
	LoginTarget:  LoginPath,
	LogoutTarget: LogoutPath,
}

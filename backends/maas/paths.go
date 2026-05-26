package maas

const (
	LoginTarget = "login"
	LoginPath   = "/MAAS/accounts/login/"
)

var TargetPaths = map[string]string{
	LoginTarget: LoginPath,
}

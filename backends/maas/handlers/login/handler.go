package login

import (
	"errors"
	"net/http"
	"net/url"

	"maas-ldap/config"
	maasldap "maas-ldap/ldap"
	"maas-ldap/proxy"
	"maas-ldap/users"
)

type loginRequest struct {
	form     url.Values
	username string
	password string
}

var (
	errDecodeRequest  = errors.New("invalid login request")
	errLDAPBind       = errors.New("ldap bind failed")
	errLDAPSearch     = errors.New("ldap search failed")
	errLDAPGroupCheck = errors.New("user is not in allowed group")
	errPasswordMap    = errors.New("user mapping not found")
	errTargetProxy    = errors.New("target proxy failed")
)

// NewHandler creates the login endpoint handler from bootstrap config.
func NewHandler(appConfig config.AppConfig, users *users.Store, target url.URL, allowedGroup string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleLogin(w, r, appConfig, users, target, allowedGroup)
	}
}

// handleLogin gates target app login behind form validation and LDAP authorization.
func handleLogin(w http.ResponseWriter, r *http.Request, appConfig config.AppConfig, users *users.Store, target url.URL, allowedGroup string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}

	login, err := decodeLoginRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := maasldap.LdapBind(login.username, login.password, appConfig.LDAP); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	allowed, err := maasldap.LdapSearch(login.username, login.password, appConfig.LDAP, allowedGroup)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if !allowed {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	mapping, ok := users.Get(login.username)

	if !ok {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Only the password is rewritten; all other form fields are preserved.
	login.form.Set("password", mapping.Secret)
	proxyBody := []byte(login.form.Encode())

	if err := proxy.ToTarget(w, r, target, proxyBody); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

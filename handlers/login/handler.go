package login

import (
	"errors"
	"net/http"
	"net/url"

	"maas-ldap/backends/maas"
	"maas-ldap/config"
	"maas-ldap/handlers/proxy"
	maasldap "maas-ldap/ldap"
	"maas-ldap/logging"
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
func NewHandler(appConfig config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleLogin(w, r, appConfig)
	}
}

// handleLogin gates target app login behind form validation and LDAP authorization.
func handleLogin(w http.ResponseWriter, r *http.Request, appConfig config.AppConfig) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}

	login, err := decodeLoginRequest(r)
	if err != nil {
		logging.Failure("-", "decode_request", errDecodeRequest)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := maasldap.LdapBind(login.username, login.password, appConfig.LDAP); err != nil {
		logging.Failure(login.username, "ldap_bind", errLDAPBind)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	allowed, err := maasldap.LdapSearch(login.username, appConfig.LDAP)
	if err != nil {
		logging.Failure(login.username, "ldap_search", errLDAPSearch)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if !allowed {
		logging.Failure(login.username, "ldap_group_check", errLDAPGroupCheck)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user, ok := appConfig.Users[login.username]

	if !ok {
		logging.Failure(login.username, "password_mapping", errPasswordMap)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Only the password is rewritten; all other form fields are preserved.
	login.form.Set("password", user.Password)
	proxyBody := []byte(login.form.Encode())

	if err := proxy.ToTarget(w, r, appConfig, maas.LoginEndpoint, proxyBody); err != nil {
		logging.Failure(login.username, "target_proxy", errTargetProxy)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

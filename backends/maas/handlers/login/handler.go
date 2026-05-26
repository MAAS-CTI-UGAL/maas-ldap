package login

import (
	"errors"
	"net/http"
	"net/url"

	maaserror "maas-ldap/backends/maas/handlers"
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
	errInvalidMethod  = errors.New("invalid HTTP method")
	errDecodeRequest  = errors.New("invalid login request")
	errLDAPBind       = errors.New("ldap bind failed")
	errLDAPSearch     = errors.New("ldap search failed")
	errLDAPGroupCheck = errors.New("user is not in allowed group")
	errPasswordMap    = errors.New("user mapping not found")
	errTargetProxy    = errors.New("target proxy failed")
)

const operationLogin = "login"

// NewHandler creates the login endpoint handler from bootstrap config.
func NewHandler(appConfig config.AppConfig, users *users.Store, target url.URL, allowedGroup string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleLogin(w, r, appConfig, users, target, allowedGroup)
	}
}

// handleLogin gates target app login behind form validation and LDAP authorization.
func handleLogin(w http.ResponseWriter, r *http.Request, appConfig config.AppConfig, users *users.Store, target url.URL, allowedGroup string) {
	if r.Method != http.MethodPost {
		maaserror.WriteError(w, operationLogin, errInvalidMethod, nil, http.StatusBadRequest)
		return
	}

	login, err := decodeLoginRequest(r)
	if err != nil {
		maaserror.WriteError(w, operationLogin, errDecodeRequest, err, http.StatusBadRequest)
		return
	}

	if err := maasldap.LdapBind(login.username, login.password, appConfig.LDAP); err != nil {
		maaserror.WriteError(w, operationLogin, errLDAPBind, err, http.StatusBadRequest)
		return
	}

	allowed, err := maasldap.LdapSearch(login.username, login.password, appConfig.LDAP, allowedGroup)
	if err != nil {
		maaserror.WriteError(w, operationLogin, errLDAPSearch, err, http.StatusBadRequest)
		return
	}

	if !allowed {
		maaserror.WriteError(w, operationLogin, errLDAPGroupCheck, nil, http.StatusBadRequest)
		return
	}

	mapping, ok := users.Get(login.username)

	if !ok {
		maaserror.WriteError(w, operationLogin, errPasswordMap, nil, http.StatusBadRequest)
		return
	}

	// Only the password is rewritten; all other form fields are preserved.
	login.form.Set("password", mapping.Secret)
	proxyBody := []byte(login.form.Encode())

	if err := proxy.ToTarget(w, r, target, proxyBody); err != nil {
		maaserror.WriteError(w, operationLogin, errTargetProxy, err, http.StatusInternalServerError)
		return
	}
}

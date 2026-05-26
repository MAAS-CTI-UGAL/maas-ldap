package login

import (
	"errors"
	"net/http"
	"net/url"

	maaserror "maas-ldap/backends/maas/handlers"
	"maas-ldap/config"
	maasldap "maas-ldap/ldap"
	"maas-ldap/proxy"
)

var (
	errInvalidMethod = errors.New("invalid HTTP method")
	errDecodeRequest = errors.New("invalid login request")
	errLDAPBind      = errors.New("ldap bind failed")
	errLDAPSearch    = errors.New("ldap search failed")
	errTargetProxy   = errors.New("target proxy failed")
)

// NewHandler creates the login endpoint handler from bootstrap config.
func NewHandler(appConfig config.AppConfig, target url.URL, allowedGroup string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleLogin(w, r, appConfig, target, allowedGroup)
	}
}

// handleLogin gates target app login behind form validation and LDAP authorization.
func handleLogin(w http.ResponseWriter, r *http.Request, appConfig config.AppConfig, target url.URL, allowedGroup string) {
	if r.Method != http.MethodPost {
		maaserror.WriteError(w, r.URL.Path, errInvalidMethod, nil, http.StatusBadRequest)
		return
	}

	form, err := decodeLoginRequest(r)
	if err != nil {
		maaserror.WriteError(w, r.URL.Path, errDecodeRequest, err, http.StatusBadRequest)
		return
	}

	username := form.Get("username")
	password := form.Get("password")

	if err := maasldap.LdapBind(username, password, appConfig.LDAP); err != nil {
		maaserror.WriteError(w, r.URL.Path, errLDAPBind, err, http.StatusBadRequest)
		return
	}

	entry, err := maasldap.LdapSearch(username, password, appConfig.LDAP, []string{"memberOf", "primaryTelexNumber"})
	if err != nil {
		maaserror.WriteError(w, r.URL.Path, errLDAPSearch, err, http.StatusBadRequest)
		return
	}

	allowed, err := checkAllowedGroup(entry, allowedGroup)
	if err != nil {
		maaserror.WriteError(w, r.URL.Path, errLDAPSearch, err, http.StatusBadRequest)
		return
	}

	if !allowed {
		maaserror.WriteError(w, r.URL.Path, errLDAPSearch, errLDAPGroupCheck, http.StatusBadRequest)
		return
	}

	maasPassword, err := maasPassword(entry)
	if err != nil {
		maaserror.WriteError(w, r.URL.Path, errLDAPSearch, err, http.StatusBadRequest)
		return
	}

	// Only the password is rewritten; all other form fields are preserved.
	form.Set("password", maasPassword)
	proxyBody := []byte(form.Encode())

	if err := proxy.ToTarget(w, r, target, proxyBody); err != nil {
		maaserror.WriteError(w, r.URL.Path, errTargetProxy, err, http.StatusInternalServerError)
		return
	}
}

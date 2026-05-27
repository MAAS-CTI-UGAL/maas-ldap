package login

import (
	"net/http"
	"net/url"

	maaserror "maas-ldap/backends/maas/handlers"
	"maas-ldap/config"
	maasldap "maas-ldap/ldap"
	"maas-ldap/proxy"
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
		w.Header().Set("Allow", http.MethodPost)
		maaserror.WriteError(w, r.URL.Path, "Invalid HTTP method", "This page only accepts login submissions.", nil, http.StatusMethodNotAllowed)
		return
	}

	form, err := decodeLoginRequest(r)
	if err != nil {
		maaserror.WriteError(w, r.URL.Path, "Invalid login request", "Please submit the login form again.", err, http.StatusBadRequest)
		return
	}

	username := form.Get("username")
	password := form.Get("password")

	if err := maasldap.LdapBind(username, password, appConfig.LDAP); err != nil {
		maaserror.WriteError(w, r.URL.Path, "LDAP bind failed", "We could not sign you in. Please check your username and password.", err, http.StatusUnauthorized)
		return
	}

	entry, err := maasldap.LdapSearch(username, password, appConfig.LDAP, []string{"memberOf", "primaryTelexNumber"})
	if err != nil {
		maaserror.WriteError(w, r.URL.Path, "LDAP search failed", "We could not verify your MAAS access. Please try again or contact an administrator.", err, http.StatusBadRequest)
		return
	}

	if !checkAllowedGroup(entry, allowedGroup) {
		maaserror.WriteError(w, r.URL.Path, "User is not in the allowed LDAP group", "You are not allowed to access MAAS.", errLDAPGroupCheck, http.StatusForbidden)
		return
	}

	maasPassword, err := maasPassword(entry)
	if err != nil {
		maaserror.WriteError(w, r.URL.Path, "LDAP entry is missing the MAAS password", "We could not verify your MAAS access. Please try again or contact an administrator.", err, http.StatusInternalServerError)
		return
	}

	// Only the password is rewritten; all other form fields are preserved.
	form.Set("password", maasPassword)
	proxyBody := []byte(form.Encode())

	if err := proxy.ToTarget(w, r, target, proxyBody); err != nil {
		maaserror.WriteError(w, r.URL.Path, "Target proxy failed", "MAAS login is temporarily unavailable. Please try again later.", err, http.StatusBadGateway)
		return
	}
}

package login

import (
	"errors"
	"net/http"
	"net/url"

	"maas-ldap/config"
	maasldap "maas-ldap/ldap"
	"maas-ldap/logging"
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
	recorder := &logging.StatusRecorder{ResponseWriter: w}
	event := logging.NewLoginEvent(r, target.String())
	defer func() {
		logging.LogLoginEvent(event, recorder.StatusCode())
	}()

	if r.Method != http.MethodPost {
		event.FailedStep = "method_check"
		event.Error = "method not allowed"
		http.Error(recorder, "Method not allowed", http.StatusBadRequest)
		return
	}

	login, err := decodeLoginRequest(r)
	if err != nil {
		event.FailedStep = "decode_request"
		event.Error = errDecodeRequest.Error()
		http.Error(recorder, "Bad request", http.StatusBadRequest)
		return
	}
	event.Username = login.username
	event.Body = redactedForm(login.form).Encode()

	if err := maasldap.LdapBind(login.username, login.password, appConfig.LDAP); err != nil {
		event.FailedStep = "ldap_bind"
		event.Error = err.Error()
		http.Error(recorder, "Bad request", http.StatusBadRequest)
		return
	}

	allowed, err := maasldap.LdapSearch(login.username, login.password, appConfig.LDAP, allowedGroup)
	if err != nil {
		event.FailedStep = "ldap_search"
		event.Error = err.Error()
		http.Error(recorder, "Bad request", http.StatusBadRequest)
		return
	}

	if !allowed {
		event.FailedStep = "ldap_group_check"
		event.Error = errLDAPGroupCheck.Error()
		http.Error(recorder, "Bad request", http.StatusBadRequest)
		return
	}

	mapping, ok := users.Get(login.username)

	if !ok {
		event.FailedStep = "username_mapping"
		event.Error = errPasswordMap.Error()
		http.Error(recorder, "Bad request", http.StatusBadRequest)
		return
	}

	// Only the password is rewritten; all other form fields are preserved.
	login.form.Set("password", mapping.Secret)
	proxyBody := []byte(login.form.Encode())

	if err := proxy.ToTarget(recorder, r, target, proxyBody); err != nil {
		event.FailedStep = "reverse_proxy"
		event.Error = errTargetProxy.Error()
		http.Error(recorder, "Internal server error", http.StatusInternalServerError)
		return
	}
	event.Outcome = "proxied"
}

func redactedForm(form url.Values) url.Values {
	redacted := url.Values{}
	for key, values := range form {
		redacted[key] = append([]string(nil), values...)
	}
	if _, ok := redacted["password"]; ok {
		redacted.Set("password", "<redacted>")
	}
	return redacted
}

package handlers

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

var (
	errUnexpectedContentType = errors.New("expected application/x-www-form-urlencoded")
	errEmptyUsername         = errors.New("username must be a non-empty string")
	errEmptyPassword         = errors.New("password must be a non-empty string")
)

// decodeLoginRequest validates and decodes the target app login form payload.
func decodeLoginRequest(r *http.Request) (loginRequest, error) {
	if !isFormContentType(r.Header.Get("Content-Type")) {
		return loginRequest{}, errUnexpectedContentType
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return loginRequest{}, err
	}

	form, err := url.ParseQuery(string(body))
	if err != nil {
		return loginRequest{}, err
	}

	username := form.Get("username")
	if username == "" {
		return loginRequest{}, errEmptyUsername
	}
	password := form.Get("password")
	if password == "" {
		return loginRequest{}, errEmptyPassword
	}

	return loginRequest{
		form:     form,
		username: username,
		password: password,
	}, nil
}

// isFormContentType accepts application/x-www-form-urlencoded, including charset parameters.
func isFormContentType(contentType string) bool {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	return strings.EqualFold(mediaType, "application/x-www-form-urlencoded")
}

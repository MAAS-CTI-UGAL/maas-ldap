package maas_manager

import (
	"encoding/json"
	"errors"
	"mime"
	"net/http"
	"strings"
)

var (
	errUnexpectedContentType = errors.New("expected application/json")
	errEmptyUsername         = errors.New("username must be a non-empty string")
	errEmptyPassword         = errors.New("password must be a non-empty string")
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type managerLoginRequest struct {
	Username string `json:"username"`
}

func decodeLoginRequest(r *http.Request) (loginRequest, error) {
	if !isJSONContentType(r.Header.Get("Content-Type")) {
		return loginRequest{}, errUnexpectedContentType
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return loginRequest{}, err
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		return loginRequest{}, errEmptyUsername
	}

	if strings.TrimSpace(req.Password) == "" {
		return loginRequest{}, errEmptyPassword
	}

	return req, nil
}

func isJSONContentType(contentType string) bool {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	return strings.EqualFold(mediaType, "application/json")
}

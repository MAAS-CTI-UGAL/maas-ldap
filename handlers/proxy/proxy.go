package proxy

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"maas-ldap/config"
)

const (
	// Context values keep the shared ReverseProxy stateless between requests.
	proxyTargetKey proxyContextKey = "proxy_target"
	proxyBodyKey   proxyContextKey = "proxy_body"
)

var errMissingTargetEndpoint = errors.New("missing target endpoint URL")

var targetProxy = &httputil.ReverseProxy{
	ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
		if recorder, ok := w.(*proxyErrorRecorder); ok {
			recorder.record(err)
		}
	},
	Rewrite: func(proxyRequest *httputil.ProxyRequest) {
		target := proxyRequest.In.Context().Value(proxyTargetKey).(url.URL)
		body := proxyRequest.In.Context().Value(proxyBodyKey).([]byte)

		// Point the outbound request at the backend while keeping the inbound query string.
		proxyRequest.SetURL(&url.URL{
			Scheme: target.Scheme,
			Host:   target.Host,
			Path:   target.Path,
		})
		proxyRequest.Out.URL.RawQuery = proxyRequest.In.URL.RawQuery
		// Replay the validated form with the mapped backend password.
		proxyRequest.Out.Method = http.MethodPost
		proxyRequest.Out.Body = io.NopCloser(bytes.NewReader(body))
		proxyRequest.Out.ContentLength = int64(len(body))
	},
}

type proxyContextKey string

type proxyErrorRecorder struct {
	http.ResponseWriter
	err error
}

// record saves proxy errors so ToTarget can return them to the caller.
func (r *proxyErrorRecorder) record(err error) {
	r.err = err
}

// ToTarget forwards the rewritten login request and streams the target response.
func ToTarget(
	w http.ResponseWriter,
	r *http.Request,
	backend config.BackendConfig,
	endpointKey string,
	body []byte,
) error {
	target, ok := backend.URLs[endpointKey]
	if !ok {
		return errMissingTargetEndpoint
	}

	ctx := context.WithValue(r.Context(), proxyTargetKey, target)
	ctx = context.WithValue(ctx, proxyBodyKey, body)
	proxyRequest := r.WithContext(ctx)

	recorder := &proxyErrorRecorder{ResponseWriter: w}
	targetProxy.ServeHTTP(recorder, proxyRequest)
	if recorder.err != nil {
		return recorder.err
	}

	return nil
}

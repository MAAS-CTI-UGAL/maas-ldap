package proxy

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ToTarget forwards the rewritten login request and streams the target response.
func ToTarget(
	w http.ResponseWriter,
	r *http.Request,
	target url.URL,
	body []byte,
) error {
	var proxyErr error

	targetProxy := &httputil.ReverseProxy{
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			proxyErr = err
		},
		Rewrite: func(proxyRequest *httputil.ProxyRequest) {
			// Point the outbound request at the backend while keeping the inbound query string.
			proxyRequest.SetURL(&url.URL{
				Scheme: target.Scheme,
				Host:   target.Host,
				Path:   target.Path,
			})
			proxyRequest.Out.URL.RawQuery = proxyRequest.In.URL.RawQuery

			// Replay the validated form with the mapped backend password.
			proxyRequest.Out.Method = proxyRequest.In.URL.RawQuery
			proxyRequest.Out.Body = io.NopCloser(bytes.NewReader(body))
			proxyRequest.Out.ContentLength = int64(len(body))
		},
	}

	targetProxy.ServeHTTP(w, r)
	return proxyErr
}

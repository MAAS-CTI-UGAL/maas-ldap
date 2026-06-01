package proxy

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ToTarget forwards a request and streams the target response.
// A nil body preserves the inbound request body; a non-nil body rewrites it.
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
			// Set the exact target URL; SetURL would append the inbound path.
			proxyRequest.Out.URL.Scheme = target.Scheme
			proxyRequest.Out.URL.Host = target.Host
			proxyRequest.Out.URL.Path = target.Path
			proxyRequest.Out.URL.RawQuery = proxyRequest.In.URL.RawQuery
			proxyRequest.Out.Host = target.Host

			proxyRequest.Out.Method = proxyRequest.In.Method
			if body == nil {
				return
			}

			// Replay the validated form with the backend credential.
			proxyRequest.Out.Body = io.NopCloser(bytes.NewReader(body))
			proxyRequest.Out.ContentLength = int64(len(body))
		},
	}

	targetProxy.ServeHTTP(w, r)
	return proxyErr
}

package handlers

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

var errMissingLoginEndpoint = errors.New("missing target login endpoint URL")

type proxyContextKey string

const (
	proxyTargetKey proxyContextKey = "proxy_target"
	proxyBodyKey   proxyContextKey = "proxy_body"
)

var targetProxy = &httputil.ReverseProxy{
	ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
		if recorder, ok := w.(*proxyErrorRecorder); ok {
			recorder.record(err)
		}
	},
	Rewrite: func(proxyRequest *httputil.ProxyRequest) {
		target := proxyRequest.In.Context().Value(proxyTargetKey).(url.URL)
		body := proxyRequest.In.Context().Value(proxyBodyKey).([]byte)

		proxyRequest.SetURL(&url.URL{
			Scheme: target.Scheme,
			Host:   target.Host,
			Path:   target.Path,
		})
		proxyRequest.Out.URL.RawQuery = proxyRequest.In.URL.RawQuery
		proxyRequest.Out.Method = http.MethodPost
		proxyRequest.Out.Body = io.NopCloser(bytes.NewReader(body))
		proxyRequest.Out.ContentLength = int64(len(body))
	},
}

type proxyErrorRecorder struct {
	http.ResponseWriter
	err error
}

func (r *proxyErrorRecorder) record(err error) {
	r.err = err
}

// proxyToTarget forwards the rewritten login request and streams the target response.
func proxyToTarget(
	w http.ResponseWriter,
	r *http.Request,
	appConfig config.AppConfig,
	body []byte,
) error {
	target, ok := appConfig.MAAS.URLs[config.EndpointLogin]
	if !ok {
		return errMissingLoginEndpoint
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

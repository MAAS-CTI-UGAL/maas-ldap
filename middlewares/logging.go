package middlewares

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &logResponseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(lrw, r)

		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, lrw.status, time.Since(start))
	})
}

type logResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *logResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

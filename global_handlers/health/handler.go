package health

import "net/http"

// NewHandler creates the health check handler.
func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleHealth(w, r)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

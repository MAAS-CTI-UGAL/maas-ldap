package health

import "net/http"

const path = "/health"

// AddRoutes registers health check routes.
func AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc(path, NewHandler())
}

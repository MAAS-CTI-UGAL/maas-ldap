package errorwriter

import (
	"log"
	"net/http"
)

// ErrorWriter logs a backend operation failure and writes the public HTTP error response.
type ErrorWriter func(http.ResponseWriter, string, string, string, error, int)

// New returns an error writer decorated with backend context.
func New(backend string) ErrorWriter {
	return func(w http.ResponseWriter, operation string, logMessage string, userMessage string, internalErr error, statusCode int) {
		if internalErr != nil {
			log.Printf("%s backend %s failed: %s: %v", backend, operation, logMessage, internalErr)
		} else {
			log.Printf("%s backend %s failed: %s", backend, operation, logMessage)
		}

		http.Error(w, userMessage, statusCode)
	}
}

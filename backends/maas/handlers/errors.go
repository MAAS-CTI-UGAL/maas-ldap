package handlers

import (
	"log"
	"net/http"
)

// WriteError logs a MAAS operation failure and writes the public HTTP error response.
// responseErr is sent to the client; internalErr is only logged so lower-level
// failure details can be diagnosed without exposing them in the response body.
func WriteError(w http.ResponseWriter, operation string, responseErr error, internalErr error, statusCode int) {
	if internalErr != nil {
		log.Printf("%s failed: %s: %v", operation, responseErr, internalErr)
	} else {
		log.Printf("%s failed: %s", operation, responseErr)
	}

	http.Error(w, responseErr.Error(), statusCode)
}

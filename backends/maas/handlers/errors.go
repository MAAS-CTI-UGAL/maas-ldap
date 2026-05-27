package handlers

import (
	"log"
	"net/http"
)

// WriteError logs a MAAS operation failure and writes the public HTTP error response.
// logMessage is only logged; userMessage is sent to the client.
func WriteError(w http.ResponseWriter, operation string, logMessage string, userMessage string, internalErr error, statusCode int) {
	if internalErr != nil {
		log.Printf("%s failed: %s: %v", operation, logMessage, internalErr)
	} else {
		log.Printf("%s failed: %s", operation, logMessage)
	}

	http.Error(w, userMessage, statusCode)
}

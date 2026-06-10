package maas

import "maas-ldap/backends/errorwriter"

// WriteError logs a MAAS operation failure and writes the public HTTP error response.
// logMessage is only logged; userMessage is sent to the client.
var WriteError = errorwriter.New("maas")

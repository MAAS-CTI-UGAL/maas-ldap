package maas_manager

import "maas-ldap/backends/errorwriter"

// WriteError logs a maas-manager operation failure and writes the public HTTP error response.
var WriteError = errorwriter.New("maas-manager")

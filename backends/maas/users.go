package maas

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"maas-ldap/users"
)

var errEmptyUsernameMapping = errors.New("maas_user_mappings contains an empty username")

// LoadUserMappings loads MAAS password mappings keyed by LDAP username.
func LoadUserMappings(database *sql.DB) ([]users.Mapping, error) {
	rows, err := database.Query("SELECT username, maas_password FROM maas_user_mappings")
	if err != nil {
		return nil, fmt.Errorf("query maas user mappings: %w", err)
	}
	defer rows.Close()

	mappings := []users.Mapping{}
	for rows.Next() {
		var username string
		var password string
		if err := rows.Scan(&username, &password); err != nil {
			return nil, fmt.Errorf("scan maas user mapping: %w", err)
		}

		if strings.TrimSpace(username) == "" {
			return nil, errEmptyUsernameMapping
		}
		if strings.TrimSpace(password) == "" {
			return nil, fmt.Errorf("maas_user_mappings row for %q has an empty maas_password", username)
		}

		mappings = append(mappings, users.Mapping{
			Username: username,
			Secret:   password,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate maas user mappings: %w", err)
	}

	return mappings, nil
}

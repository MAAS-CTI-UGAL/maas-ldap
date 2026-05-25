# Environment Configuration

The app loads configuration from a `.env` file at startup. If `.env` cannot be
loaded, or if any required value is missing, startup stops with a fatal error.

## Example

```env
LDAP_URL=ldap://ldap.example.internal:389
LDAP_UPN_SUFFIX=example.internal
LDAP_BASE_DN=DC=example,DC=internal
MAAS_URL=https://maas.example.internal
MAAS_LDAP_ALLOWED_GROUP=MaaS_Allowed
DB_PATH=/var/lib/maas-ldap/maas-ldap.db

PORT=8080
```

## MAAS CTI Example

For the CTI MAAS deployment, these values are known:

```env
LDAP_URL=ldap://10.13.11.1:389
LDAP_UPN_SUFFIX=cti.ugal.ro
LDAP_BASE_DN=DC=CTI,DC=UGAL,DC=RO
MAAS_URL=http://10.13.201.10:5240
MAAS_LDAP_ALLOWED_GROUP=MaaS_Allowed
```

Local development can use repository-local runtime files:

```env
DB_PATH=./maas-ldap.db
LOG_PATH=./maas-ldap.log
```

Production should normally use persistent system paths, for example:

```env
DB_PATH=/var/lib/maas-ldap/maas-ldap.db
```

When running as a `systemd` service, omit `LOG_PATH` in production so logs go
to stderr and are captured by `journald`. Use `journalctl -u maas-ldap` or
`journalctl -u maas-ldap -f` to inspect them.

## Required Values

The global LDAP values configure how the app connects to LDAP and searches for
users. Backend-specific LDAP values configure which LDAP group authorizes access
to each backend.

`LDAP_URL`

LDAP server URL, including scheme and port.

Examples:

```env
LDAP_URL=ldap://ldap.example.internal:389
LDAP_URL=ldaps://ldap.example.internal:636
```

`LDAP_UPN_SUFFIX`

Suffix used when binding to LDAP. The app binds as:

```text
<submitted-username>@<LDAP_UPN_SUFFIX>
```

Example:

```env
LDAP_UPN_SUFFIX=example.internal
```

`LDAP_BASE_DN`

LDAP base DN used when searching for the submitted user.

Example:

```env
LDAP_BASE_DN=DC=example,DC=internal
```

`<BACKEND_NAME>_LDAP_ALLOWED_GROUP`

LDAP group required for backend access. Each backend owns its own allowed group
using this naming pattern. For the current MAAS backend, the required variable
is `MAAS_LDAP_ALLOWED_GROUP`.

Examples:

```env
MAAS_LDAP_ALLOWED_GROUP=MaaS_Allowed
MAAS_LDAP_ALLOWED_GROUP=CN=MaaS_Allowed,OU=Groups,DC=example,DC=internal
GRAFANA_LDAP_ALLOWED_GROUP=Grafana_Allowed
```

The value can be either a short group CN or a full group DN. Use a full DN when
different OUs may contain groups with the same CN.

`MAAS_URL`

Base URL for the MAAS backend. Include the scheme and host, but do not include
the login path. The app appends `/MAAS/accounts/login/` itself.

Examples:

```env
MAAS_URL=https://maas.example.internal
MAAS_URL=http://maas.example.internal:5240
```

`DB_PATH`

SQLite database file path used for MAAS user password mappings.

Example:

```env
DB_PATH=/var/lib/maas-ldap/maas-ldap.db
```

The value must point to a database file, not only a directory. For a database in
the current working directory, use a filename such as:

```env
DB_PATH=./maas-ldap.db
```

If the database file does not exist, SQLite creates it. The parent directory
must already exist and be writable by the user running the app.

## Optional Values

`PORT`

Listen port or address. If unset, the app listens on `0.0.0.0:9090`.

Examples:

```env
PORT=8080
PORT=:8080
PORT=127.0.0.1:8080
```

`LOG_PATH`

Optional path to a log file for local development or deployments that
deliberately want file logging. If unset, logs are written only to stderr.
For production `systemd` services, prefer leaving this unset and reading logs
from the journal.

Example:

```env
LOG_PATH=./maas-ldap.log
```

When set, logs are written to both stderr and the configured file. The parent
directory must exist, and the user running the app must be able to append to the
file.

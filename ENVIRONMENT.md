# Environment Configuration

The app loads configuration from a `.env` file at startup. If `.env` cannot be
loaded, or if any required value is missing, startup stops with a fatal error.

## Example

```env
LDAP_URL=ldap://ldap.example.internal:389
LDAP_UPN_SUFFIX=example.internal
LDAP_BASE_DN=DC=example,DC=internal
LDAP_ALLOWED_GROUP=MaaS_Allowed
MAAS_URL=https://maas.example.internal
DB_PATH=/var/lib/maas-ldap/maas-ldap.db

PORT=8080
LOG_PATH=/var/log/maas-ldap/maas-ldap.log
```

## Required Values

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

`LDAP_ALLOWED_GROUP`

LDAP group required for access. This can be either a short group CN or a full
group DN.

Examples:

```env
LDAP_ALLOWED_GROUP=MaaS_Allowed
LDAP_ALLOWED_GROUP=CN=MaaS_Allowed,OU=Groups,DC=example,DC=internal
```

Use a full DN when different OUs may contain groups with the same CN.

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

Path to a log file. If unset, logs are written only to stderr.

Example:

```env
LOG_PATH=/var/log/maas-ldap/maas-ldap.log
```

When set, logs are written to both stderr and the configured file. The parent
directory must exist, and the user running the app must be able to append to the
file.

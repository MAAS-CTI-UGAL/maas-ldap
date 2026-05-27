# maas-ldap

LDAP-gated login proxy for MAAS.

The current implementation supports MAAS. The service exposes the MAAS login
endpoint, validates submitted credentials against LDAP, checks group
membership, rewrites the submitted password to the MAAS password stored in the
LDAP `primaryTelexNumber` attribute, and proxies the request to the real MAAS
backend.

The project is structured so additional login backends can be added to the
backend registry.

## Architecture

Global configuration owns shared LDAP connection and search settings. The
backend registry owns each backend's login path, target URL environment key,
allowed LDAP group environment key, and login handler.

`BACKENDS` controls which registered backends are active. Shared backend helpers
build full target URLs from a backend base URL and the login path defined by
that backend.

## Login Flow

`POST /MAAS/accounts/login/`

1. Accepts `application/x-www-form-urlencoded` login requests.
2. Reads `username` and `password` from the request body.
3. Binds to LDAP as `username@LDAP_UPN_SUFFIX` with the submitted password.
4. Binds again with the submitted credentials, then searches LDAP under
   `LDAP_BASE_DN` for:

   ```text
   (&(objectClass=user)(sAMAccountName=<username>))
   ```

5. Requires exactly one user result.
6. Requires one `memberOf` value to match `MAAS_LDAP_ALLOWED_GROUP`.
   `MAAS_LDAP_ALLOWED_GROUP` can be a full group DN or a short group CN.
7. Requires exactly one non-empty `primaryTelexNumber` value on the LDAP user.
8. Replaces only the `password` form value with that `primaryTelexNumber`
   value.
9. Proxies the request to `${MAAS_URL}/MAAS/accounts/login/`.
10. Streams the MAAS response back to the client.

Validation, LDAP, group, and LDAP MAAS-password failures return
`400 Bad Request`.
Unexpected proxy failures return `500 Internal Server Error`.

## Configuration

Configuration is loaded from `.env` at startup.

Required:

```env
BACKENDS=maas

LDAP_URL=ldap://example.internal:389
LDAP_UPN_SUFFIX=example.internal
LDAP_BASE_DN=DC=example,DC=internal

MAAS_URL=https://maas.example.internal
MAAS_LDAP_ALLOWED_GROUP=MaaS_Allowed
```

`BACKENDS` is a comma-separated list of registered backend names to enable.
The current registered backend is `maas`.

`LDAP_URL`, `LDAP_UPN_SUFFIX`, and `LDAP_BASE_DN` configure the global LDAP
server and user search shared by all backends. Backend-specific authorization
groups use the `<BACKEND_NAME>_LDAP_ALLOWED_GROUP` pattern; the current MAAS
backend requires `MAAS_LDAP_ALLOWED_GROUP`.

Allowed groups can be a short CN or a full DN. For example:

```env
MAAS_LDAP_ALLOWED_GROUP=MaaS_Allowed
MAAS_LDAP_ALLOWED_GROUP=CN=MaaS_Allowed,OU=Groups,DC=example,DC=internal
```

Use a full DN when the deployment needs to distinguish between groups with the
same CN in different OUs.

The LDAP user entry must also contain exactly one non-empty
`primaryTelexNumber` value. The proxy treats this value as the MAAS password and
does not log it.

Optional:

```env
PORT=8080
```

If `LOG_PATH` is not set, logs are written only to stderr. This is the
recommended production setup for `systemd` services because stderr is captured
by `journald`.

## Run

```bash
go run .
```

Build a binary:

```bash
go build -o maas-ldap .
```

By default, the service listens on:

```text
:9090
```

## Logging

The app always writes logs to stderr. In production under `systemd`, leave
`LOG_PATH` unset and read logs from `journald`.

When running under `systemd`, stderr is captured by the journal:

```bash
journalctl -u maas-ldap
journalctl -u maas-ldap -f
```

For local development or deployments that deliberately want a simple file log,
set `LOG_PATH`. When set, the same log lines are also appended to that file.

HTTP requests emit one access log line:

```text
POST /MAAS/accounts/login/ -> 204 (12.345ms)
```

## Service User Notes

If `LOG_PATH` is deliberately enabled, create its parent directory and make it
writable by the service user. Production `systemd` deployments should normally
omit `LOG_PATH` and use `journalctl` instead.

## Check

```bash
go test ./...
```

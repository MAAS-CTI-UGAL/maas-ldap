# maas-ldap Architecture

`maas-ldap` is a Go HTTP service that gates backend login requests through
LDAP authentication and group authorization before proxying them to MAAS or
related services. It does not own user sessions for the target applications.
Instead, it validates the submitted credentials, rewrites the target login
payload when needed, and streams the target response back to the caller.

This document describes the repository structure and operating model. It is a
maintainer guide, not a status tracker.

## Startup

`main.go` is intentionally small:

- load app-wide startup configuration through `config.Bootstrap()`
- load enabled backend configuration through `backends.LoadEnabledConfigs()`
- create an `http.ServeMux`
- register global and backend routes with `AddRoutes(...)`
- wrap the mux with request logging middleware
- serve HTTP on `AppConfig.ListenAddress`

The listen address is loaded from `PORT` and defaults to `127.0.0.1:9090`, so
the service is reachable only from the local host. A `PORT` value may override
this address when a different deployment topology explicitly requires it.

Configuration validation happens during startup. App-wide LDAP and logging
settings are loaded by `config.Bootstrap()`. Backend target URLs and allowed
LDAP groups are loaded separately by `backends.LoadEnabledConfigs()` for each
backend listed in `BACKENDS`.

The app currently calls `godotenv.Load()` during bootstrap and stops startup if
the `.env` file cannot be loaded.

## Package Layout

The repository is organized around request flow boundaries:

```text
backends/                 backend registry, config loading, and backend route registration
backends/errorwriter/     shared safe handler error writer
backends/maas/            MAAS form login validation, LDAP authorization, and password rewrite
backends/maas_manager/    maas-manager JSON login validation and username-only proxying
config/                   .env loading, app-wide environment validation, LDAP config, logging
global_handlers/          global endpoints such as health checks
ldap/                     LDAP bind, search, user filter, and group membership checks
middlewares/              HTTP middleware
proxy/                    shared reverse proxy forwarding
```

Keep backend-specific request parsing and login behavior under
`backends/<name>/`. Keep shared LDAP behavior in `ldap/` and shared forwarding
behavior in `proxy/`. Route registration lives in `routes.go` files so startup
wiring stays easy to follow.

## Routing

Top-level route registration happens in `routes.go`:

- `global_handlers.AddRoutes(mux)` registers global routes.
- `backends.AddRoutes(mux, appConfig, backendConfigs)` registers one login
  route for each enabled backend.

Global routes:

```text
GET /health
```

`/health` returns `200 OK` with a plain text `ok` body. The logging middleware
does not log health check requests.

Backend routes are enabled by `BACKENDS`:

```text
POST /MAAS/accounts/login/   enabled by BACKENDS=maas
POST /manager/api/login      enabled by BACKENDS=maas-manager
```

Each backend route also rejects non-`POST` requests with `405 Method Not
Allowed` and an `Allow: POST` header.

## Configuration

Required app-wide environment:

```text
BACKENDS
LDAP_URL
LDAP_UPN_SUFFIX
LDAP_BASE_DN
```

`BACKENDS` is a comma-separated list of registered backend names. Names are
matched case-insensitively. Empty values, duplicate backend names, and unknown
backend names fail startup.

Current backend definitions:

```text
maas          MAAS_URL           MAAS_LDAP_ALLOWED_GROUP           /MAAS/accounts/login/
maas-manager  MAAS_MANAGER_URL   MAAS_MANAGER_LDAP_ALLOWED_GROUP   /manager/api/login
```

The `maas-manager` allowed-group environment variable above reflects the current
code in `backends/registry.go`. Some older docs may still say the manager
backend shares `MAAS_LDAP_ALLOWED_GROUP`; update those docs separately if the
deployment model should remain split by backend.

Optional app-wide environment:

```text
PORT=127.0.0.1:9090
LOG_PATH
```

When `LOG_PATH` is set, logs are written to both stderr and the configured file.
Otherwise logs go to stderr. Production `systemd` deployments should normally
leave `LOG_PATH` unset and read logs from `journald`.

Backend target URL values must include scheme and host. The backend config
loader appends the backend login path while preserving any base path on the
configured URL.

## LDAP Authentication and Authorization

LDAP settings are app-wide and shared by all backends:

```text
LDAP_URL
LDAP_UPN_SUFFIX
LDAP_BASE_DN
```

The app binds to LDAP using the submitted credentials as:

```text
<submitted-username>@<LDAP_UPN_SUFFIX>
```

`ldap.LdapSearch(...)` binds with those credentials, searches under
`LDAP_BASE_DN`, and expects exactly one matching user entry. The default user
filter is:

```text
(&(objectClass=user)(sAMAccountName=<escaped username>))
```

Backend handlers choose which LDAP attributes they need for their flow. All
current backends request `memberOf`; the MAAS backend also requests
`primaryTelexNumber`.

Group authorization is handled by `ldap.CheckAllowedGroup(...)`. LDAP
`memberOf` values are full DNs. The configured allowed group can be either:

- a full DN, matched case-insensitively against `memberOf`
- a short CN, matched against the start of a `memberOf` DN

When the configured allowed group contains `=`, it is treated as a full DN and
short-CN matching is disabled.

## Backend Login Flows

The MAAS backend accepts `application/x-www-form-urlencoded` login requests. It
requires non-empty `username` and `password` fields.

For normal users, the MAAS flow is:

- search LDAP using the submitted username and password
- require membership in `MAAS_LDAP_ALLOWED_GROUP`
- require exactly one non-empty `primaryTelexNumber` value
- replace only the submitted form `password` with `primaryTelexNumber`
- proxy the rewritten form to the configured MAAS login URL

All other submitted form fields are preserved.

The MAAS username `maas_admin` bypasses LDAP checks and is proxied unchanged.
This exists because that account is not expected to have a valid LDAP account.

The `maas-manager` backend accepts `application/json` login requests:

```json
{
  "username": "alice",
  "password": "submitted-ldap-password"
}
```

It trims and validates `username`, requires a non-empty password, searches LDAP
using the submitted credentials, and requires membership in
`MAAS_MANAGER_LDAP_ALLOWED_GROUP`. It does not read or forward the submitted
password. After authorization, it proxies this JSON body to `maas-manager`:

```json
{
  "username": "alice"
}
```

`maas-manager` owns its own browser session creation after receiving the
username-only login request.

## Reverse Proxy Behavior

`proxy.ToTarget(...)` is the shared forwarding path for backend handlers.

The proxy:

- sets the exact target scheme, host, and path from backend configuration
- preserves the inbound request query string
- sets the outbound host to the target host
- preserves the inbound method
- preserves the inbound request body when called with a `nil` body
- rewrites the outbound request body and content length when called with a
  non-`nil` body

The proxy deliberately sets the target URL directly instead of using
`ReverseProxy.SetURL`, because backend login routes should forward to their
configured target path rather than appending the inbound path.

## Handler Error Policy

Backend handlers return safe public error messages and log internal details
through `backends/errorwriter`.

Use status codes consistently:

```text
400 invalid request body, content type, or required fields
401 LDAP credential/search failure
403 authenticated user is not in the allowed LDAP group
405 invalid method, with Allow header
500 local request construction or LDAP entry data problem
502 backend proxy failure
```

Internal log messages may include backend names, route paths, and error causes.
They must not include submitted passwords, LDAP secrets, MAAS passwords from
`primaryTelexNumber`, rewritten proxy bodies, API keys, or other credentials.

## Operational Notes

Serve `maas-ldap` behind TLS whenever browser credentials cross an untrusted
network. The service receives user passwords before LDAP validation, and the
MAAS backend may substitute a backend-specific password from LDAP.

Keep `.env` readable only by trusted deployment users. It contains LDAP
connection information, backend target URLs, and authorization group names.

Production deployments should normally rely on stderr logging captured by
`journald`. The request logging middleware records method, path, status, and
duration for non-health requests.

The service is stateless with respect to backend sessions. Availability depends
on LDAP and on the enabled backend targets. Backend application cookies and
session semantics are owned by the target applications, not by `maas-ldap`.

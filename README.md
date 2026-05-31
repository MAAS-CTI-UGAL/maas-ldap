# maas-ldap

LDAP-gated login proxy for MAAS.

The service accepts MAAS login requests, validates submitted credentials
against LDAP, checks group membership, replaces the submitted password with the
MAAS password stored in LDAP `primaryTelexNumber`, and proxies the login request
to the real MAAS backend.

## Documentation

- [docs/PREREQUISITES.md](docs/PREREQUISITES.md): local, runtime, and production host prerequisites
- [docs/ENVIRONMENT.md](docs/ENVIRONMENT.md): environment variables and configuration model
- [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md): GitLab CI deployment, sudoers, and operations
- [docs/GIT.md](docs/GIT.md): Git remote setup

## Configuration

Startup configuration is split between app-wide settings and enabled backend
settings. `config.Bootstrap()` loads the app config: listen address, LDAP
connection settings, and logging. `backends.LoadEnabledConfigs()` loads one
validated config object for each backend listed in `BACKENDS`.

For the current MAAS backend, `MAAS_URL` and `MAAS_LDAP_ALLOWED_GROUP` are
backend settings. Shared LDAP settings such as `LDAP_URL`, `LDAP_UPN_SUFFIX`,
and `LDAP_BASE_DN` remain app-wide settings.

## Run

```bash
go run .
```

By default, the service listens on:

```text
:9090
```

## Build

```bash
go build -o maas-ldap .
```

## Check

```bash
go test ./...
```

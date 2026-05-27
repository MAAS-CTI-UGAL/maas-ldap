# maas-ldap

LDAP-gated login proxy for MAAS.

The service accepts MAAS login requests, validates submitted credentials
against LDAP, checks group membership, replaces the submitted password with the
MAAS password stored in LDAP `primaryTelexNumber`, and proxies the login request
to the real MAAS backend.

## Documentation

- [PREREQUISITES.md](PREREQUISITES.md): local, runtime, and production host prerequisites
- [DEPLOYMENT.md](DEPLOYMENT.md): GitLab CI deployment, sudoers, and operations
- [GIT.md](GIT.md): Git remote setup

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

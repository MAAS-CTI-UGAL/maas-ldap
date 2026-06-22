# maas-ldap Prerequisites

This document covers the prerequisites needed to run, configure, and deploy
`maas-ldap`.

## Local Development

Install:

- Go
- Git

The app is configured from a `.env` file in the repository root.

## Runtime Configuration

Required `.env` values:

```env
BACKENDS=maas,maas-manager

LDAP_URL=ldap://example.internal:389
LDAP_UPN_SUFFIX=example.internal
LDAP_BASE_DN=DC=example,DC=internal

MAAS_URL=https://maas.example.internal
MAAS_MANAGER_URL=http://127.0.0.1:9091
MAAS_LDAP_ALLOWED_GROUP=MaaS_Allowed
```

`BACKENDS` is a comma-separated list of registered backend names to enable.
The current registered backends are `maas` and `maas-manager`.

`LDAP_URL`, `LDAP_UPN_SUFFIX`, and `LDAP_BASE_DN` configure the app-wide LDAP
server and user search shared by all backends. Backend-specific values are
loaded for each backend listed in `BACKENDS`; the current MAAS backend uses
`MAAS_URL`, and the `maas-manager` backend uses `MAAS_MANAGER_URL`. Both
backends currently use `MAAS_LDAP_ALLOWED_GROUP`.

Allowed groups can be a short CN or a full DN:

```env
MAAS_LDAP_ALLOWED_GROUP=MaaS_Allowed
MAAS_LDAP_ALLOWED_GROUP=CN=MaaS_Allowed,OU=Groups,DC=example,DC=internal
```

Use a full DN when the deployment needs to distinguish between groups with the
same CN in different OUs.

The LDAP user entry must contain exactly one non-empty `primaryTelexNumber`
value. The proxy treats this value as the MAAS password and does not log it.
The `maas-manager` backend does not use `primaryTelexNumber`; it validates LDAP
credentials and group membership, then forwards only the username to
`maas-manager`.

Optional values:

```env
PORT=8080
```

If `LOG_PATH` is not set, logs are written only to stderr. This is the
recommended production setup for `systemd` services because stderr is captured
by `journald`.

## Production Host

Production deployment uses a shell GitLab Runner on the host that runs the
service. Install Go and Git on the host before running the pipeline.

The service runs as the locked `maas` system user with group `deploy`.
One-time host setup must create the users and groups:

```bash
sudo groupadd --system deploy
sudo groupadd --system maas
sudo useradd --system --gid maas --groups deploy --home-dir /nonexistent --shell /usr/sbin/nologin --comment "MAAS service account" maas
sudo usermod -aG deploy gitlab-runner
```

If the `maas` user already exists, only ensure both users are in `deploy`:

```bash
sudo usermod -aG deploy maas
sudo usermod -aG deploy gitlab-runner
```

The exact UID and GID do not need to be hardcoded. A valid service account can
look like:

```text
maas:x:999:988:MAAS service account:/nonexistent:/usr/sbin/nologin
```

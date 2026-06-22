# maas-ldap Deployment

Production deployment uses a shell GitLab Runner on the host that runs the
service. The runner builds in its temporary checkout, installs runtime files
into `/opt/maas-ldap`, installs the systemd unit, and restarts the service.

## Host Prerequisites

See [PREREQUISITES.md](PREREQUISITES.md) for local, runtime, and production
host prerequisites. The CI job creates and repairs `/opt/maas-ldap` during
deployment.

## GitLab Variable

The GitLab CI/CD variable for the production `.env` file must be:

- name: `ENV_FILE`
- type: `File`
- protected
- not masked, because multiline file variables cannot be masked
- variable reference expansion disabled

Use these contents for the current CTI MAAS deployment:

```env
BACKENDS=maas,maas-manager
LDAP_URL=ldap://10.13.11.1:389
LDAP_UPN_SUFFIX=cti.ugal.ro
LDAP_BASE_DN=DC=CTI,DC=UGAL,DC=RO
MAAS_URL=http://10.13.201.10:5240
MAAS_MANAGER_URL=http://127.0.0.1:9091
PORT=127.0.0.1:9090
MAAS_LDAP_ALLOWED_GROUP=MaaS_Allowed
```

## Sudoers

Add this sudoers entry with `visudo`, preferably as
`/etc/sudoers.d/gitlab-runner-maas-ldap`:

```sudoers
gitlab-runner ALL=(root) NOPASSWD: /usr/bin/install -d -o maas -g deploy -m 2770 /opt/maas-ldap, /usr/bin/chown maas\:deploy /opt/maas-ldap /opt/maas-ldap/maas-ldap /opt/maas-ldap/.env /opt/maas-ldap/maas-ldap.service, /usr/bin/chmod 2770 /opt/maas-ldap, /usr/bin/chmod 0750 /opt/maas-ldap/maas-ldap, /usr/bin/chmod 0640 /opt/maas-ldap/.env, /usr/bin/chmod 0644 /opt/maas-ldap/maas-ldap.service, /usr/bin/install -m 0644 /opt/maas-ldap/maas-ldap.service /etc/systemd/system/maas-ldap.service, /usr/bin/systemctl daemon-reload, /usr/bin/systemctl enable maas-ldap.service, /usr/bin/systemctl restart maas-ldap.service, /usr/bin/systemctl status maas-ldap.service --no-pager
```

The CI job uses explicit `install`, `chown`, and `chmod` modes for deployed
files. A `umask` is not required; the explicit modes are clearer and the setgid
bit on `/opt/maas-ldap` keeps new files in the `deploy` group.

## CI Deployment Flow

The deploy job:

1. Confirms the `ENV_FILE` GitLab file variable is readable.
2. Downloads Go module dependencies.
3. Builds the `maas-ldap` binary.
4. Creates or repairs `/opt/maas-ldap`.
5. Installs the binary, `.env`, and systemd unit.
6. Applies the expected ownership and permissions.
7. Installs the unit into `/etc/systemd/system`.
8. Reloads systemd, enables the service, restarts it, and prints service status.

## Logs

Production deployments should normally leave `LOG_PATH` unset and use
`journald`:

```bash
journalctl -u maas-ldap
journalctl -u maas-ldap -f
```

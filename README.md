# maas-ldap

LDAP-gated login proxy for MAAS.

The service exposes the MAAS login endpoint, validates submitted credentials
against LDAP, checks group membership, rewrites the submitted password to the
mapped MAAS password from SQLite, and proxies the request to the real
MAAS backend.

## Login Flow

`POST /MAAS/accounts/login/`

1. Accepts `application/x-www-form-urlencoded` login requests.
2. Reads `username` and `password` from the request body.
3. Binds to LDAP as `username@LDAP_UPN_SUFFIX` with the submitted password.
4. Searches LDAP under `LDAP_BASE_DN` for:

   ```text
   (&(objectClass=user)(sAMAccountName=<username>))
   ```

5. Requires exactly one user result.
6. Requires one `memberOf` value to match `LDAP_ALLOWED_GROUP`.
   `LDAP_ALLOWED_GROUP` can be a full group DN or a short group CN.
7. Looks up the username in `maas_user_mappings`.
8. Replaces only the `password` form value with the mapped MAAS password.
9. Proxies the request to `${MAAS_URL}/MAAS/accounts/login/`.
10. Streams the MAAS response back to the client.

Validation, LDAP, group, and mapping failures return `400 Bad Request`.
Unexpected proxy failures return `500 Internal Server Error`.

## Configuration

Configuration is loaded from `.env` at startup.

Required:

```env
LDAP_URL=ldap://example.internal:389
LDAP_UPN_SUFFIX=example.internal
LDAP_BASE_DN=DC=example,DC=internal
LDAP_ALLOWED_GROUP=MaaS_Allowed
MAAS_URL=https://maas.example.internal
DB_PATH=/var/lib/maas-ldap/maas-ldap.db
```

`LDAP_ALLOWED_GROUP` also accepts a full group DN, such as
`CN=MaaS_Allowed,OU=Groups,DC=example,DC=internal`, when the deployment needs
to distinguish between groups with the same CN in different OUs.

Optional:

```env
PORT=8080
LOG_PATH=/var/log/maas-ldap/maas-ldap.log
```

If `LOG_PATH` is not set, logs are written only to stderr.

## SQLite User Mappings

The app opens `DB_PATH`, runs embedded migrations, and loads
`maas_user_mappings` into memory at startup. If the DB file does not exist,
SQLite creates it. The parent directory must already exist and be writable.

`maas_user_mappings` maps LDAP usernames to MAAS passwords:

```sql
INSERT INTO maas_user_mappings (username, maas_password)
VALUES ('some.username', 'maas-password')
ON CONFLICT(username) DO UPDATE SET
    maas_password = excluded.maas_password;
```

The `username` value must match the submitted LDAP username. Password mappings
are loaded once at startup; restart the service after updating the database.

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

The app always writes logs to stderr. If `LOG_PATH` is set, the same log lines
are also appended to that file.

When running under `systemd`, stderr is captured by the journal:

```bash
journalctl -u maas-ldap
```

Login failures use this format:

```text
user=<username> failed_step=<step> error=<error>
```

Known failure steps:

```text
decode_request
ldap_bind
ldap_search
ldap_group_check
password_mapping
target_proxy
```

## Service User Notes

If `LOG_PATH` is set and the binary runs as a dedicated `maas-ldap` user, make
sure that user can append to the configured log file. One option is to create
the log file with a group that the service user belongs to, then give the group
write permission.

Example:

```bash
sudo mkdir -p /var/log/maas-ldap
sudo touch /var/log/maas-ldap/maas-ldap.log
sudo chown root:syslog /var/log/maas-ldap/maas-ldap.log
sudo chmod 0660 /var/log/maas-ldap/maas-ldap.log
sudo usermod -aG syslog maas-ldap
```

## Check

```bash
go test ./...
```

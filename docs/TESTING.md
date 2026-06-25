# Testing

This project uses Go's standard test runner. Test files live next to the code
they cover and use the `*_test.go` suffix.

## Run The Suite

Run all tests from the repository root:

```bash
go test ./...
```

`go build` does not run tests. Use both commands when you want to verify tests
and produce a binary:

```bash
go test ./...
go build -o maas-ldap .
```

## Coverage

Generate a statement coverage profile:

```bash
go test ./... -coverprofile=coverage.out
```

Inspect function-level coverage:

```bash
go tool cover -func=coverage.out
```

Open the HTML coverage report:

```bash
go tool cover -html=coverage.out
```

## Test Scope

The suite focuses on deterministic behavior:

- backend configuration loading and validation
- enabled backend registry parsing
- request decoding for the MAAS and `maas-manager` login handlers
- LDAP filter, group matching, and local validation paths
- HTTP handler early returns and error responses
- route registration, health checks, logging middleware, and reverse proxy behavior

Tests must not call real LDAP servers or production backend services. Use
package-local helpers, `httptest`, temporary files, and environment variables
set with `t.Setenv`.

## Local HTTP Tests

Most handler tests call handlers directly with `httptest.NewRequest` and
`httptest.NewRecorder`; they do not start a server.

The reverse proxy tests use `httptest.NewServer` to create a local loopback
target. In restricted sandboxes, local socket creation may be blocked. On a
normal development machine, `go test ./...` should run without extra setup.

## Naming Conventions

Use behavior-oriented test names:

```go
func TestLoadBackendConfigBuildsTargetAndTrimsAllowedGroup(t *testing.T)
```

Use table-driven tests when several inputs exercise the same behavior, and use
`t.Run` with short scenario names. Use `t.Helper` for shared assertion helpers
so failures point to the caller.

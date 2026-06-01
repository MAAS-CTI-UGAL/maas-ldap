# Repository Instructions

- Do not write tests unless the user explicitly asks for them.
- Prefer direct string literals over constants unless reuse or clarity clearly justifies a constant.

## Project Structure

- Keep `main.go` focused on startup: load config, load backend configs, wire routes, wrap middleware, and start the HTTP server.
- Keep `.env` loading, startup validation, log setup, and typed runtime configuration in `config/`.
- Keep route registration in `routes.go` files. Top-level routes should compose global handlers and backend routes.
- Keep backend-specific login behavior under `backends/<name>/`.
- Keep shared LDAP behavior in `ldap/` and shared reverse-proxy behavior in `proxy/`.
- Avoid mixing config loading, route registration, request parsing, LDAP calls, and proxy forwarding in one file.

## Handler Style

- Expose handler constructors as `NewHandler(...) http.HandlerFunc`.
- Keep constructor closures small and delegate real work to private `handle...` functions.
- Validate HTTP method, content type, and required request input early.
- Move non-trivial request decoding and parsing into focused helper files.
- Pass config and runtime dependencies through constructors and function parameters instead of mutable package globals.

## Error Handling

- Startup and configuration failures should produce clear errors and fail at startup boundaries with `log.Fatal`.
- Use package-level `errors.New` values for reusable validation or domain errors.
- Wrap errors with `%w` when preserving the original cause matters to callers.
- Handler responses should use safe public messages and log internal details separately, following the existing `WriteError` pattern.
- Never log submitted passwords, LDAP secrets, backend credentials, or rewritten proxy credentials.

## Go Style

- Run `gofmt` on edited Go files.
- Prefer the standard library unless a dependency is already present or clearly justified.
- Keep comments short and useful. Comment exported symbols and non-obvious behavior.
- Keep functions focused on one responsibility and extract helpers when validation, parsing, or transformation logic grows.

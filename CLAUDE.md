# go-logging — Claude context

## What this library is

A small slog-backed structured logging library for Go. The public surface is intentionally narrow: a single `Logger` interface, a single `Config` struct, a single `New` constructor, plus context helpers for trace/request/correlation IDs and a logger-in-context idiom. See [README.md](README.md) for the canonical usage.

The library deliberately does not implement: handler builders, middleware chains, async log shipping, multiple formatters, redaction, YAML config, DI providers, or HTTP middleware. When those concerns arise, callers wrap their own `slog.Handler` and pass it via `Config.Handler`, or layer functionality in their application code.

## Architecture

- `Logger` interface (10 methods) — defined in `pkg/logging/logger.go`. Implementations must be safe for concurrent use.
- `slogLogger` is the only built-in implementation. It wraps `*slog.Logger`.
- `Config` is a flat struct with seven fields. `New(Config)` is the only constructor.
- Context keys are private; package-level helpers (`WithTraceID`, `TraceIDFromContext`, …) are the only access path. The `loggerKey` is what `WithContext`/`FromContext` use to stash and retrieve a logger.
- `default.go` holds the package-default logger seam (`SetDefaultLogger` + package-level `Debug`/`Info`/`Warn`/`Error`). The default is a fallback for tests and one-off scripts, not the recommended call site — production code should construct a `Logger` and inject it.

## Conventions

- **Commit messages**: Conventional Commits with optional scope. Useful scopes: `logging`, `config`, `context`, `default`, `docs`, `ci`, `test`. Breaking changes use `!` after the type.
- **Tests**: external test package (`package logging_test`) so we exercise the public surface, not internals. Use `t.Parallel()` everywhere except in `default_test.go` (mutates package-level state).
- **No new dependencies** without a strong reason. The module currently has zero external dependencies and that is a feature.
- **No global mutable state** beyond the package-default logger seam. If a future feature needs configuration, prefer passing it through `Config` over package-level setters.

## What not to add

- New levels (`Trace`, `Critical`, etc.) — slog's four are the canonical set.
- Fluent builders — `With(args ...any)` covers field attachment idiomatically.
- Multiple constructors — if a constructor with different defaults is wanted, add a `Config{...}` builder helper or document the literal at the call site.
- Async writes — wrap `Config.Output` with an async `io.Writer` at the call site if needed.
- Redaction patterns — wrap `Config.Handler` with a redacting `slog.Handler` at the call site.

## Development

```bash
task test          # go test ./... -race -count=1
task lint          # golangci-lint run ./...
task tidy          # go mod tidy && go mod verify
```

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.0.0] - 2026-06-13

v2 is a near-complete rewrite. The public API was reduced from ~150 exported
symbols to ~25, the production code from ~11,800 lines to ~330. The library now
sits as a thin layer over `log/slog` with no external dependencies. See
[README.md](README.md) for the canonical usage.

### Breaking

- **Module path renamed**: `github.com/ocrosby/go-logging` → `github.com/jedi-knights/go-logging`.
- **Go version**: minimum bumped from 1.24 to 1.26.
- **`Logger` interface reduced to 10 methods**: `Debug`/`Info`/`Warn`/`Error` plus the `*Context` variants, `With(args ...any) Logger`, and `Enabled(ctx, slog.Level) bool`. Removed: `Trace`, `Critical`, `Log`, `LogContext`, `TraceContext`, `CriticalContext`, `WithField`, `WithFields`, `Fluent`, `IsLevelEnabled`, `SetLevel`, `GetLevel`.
- **`Config` collapsed to a flat struct**: `Level`, `Format`, `Output`, `ServiceName`, `Environment`, `StaticFields`, `Handler`. The deprecated `Config`, the nested `LoggerConfig`/`CoreConfig`/`FormatterConfig`/`OutputConfig`, and all their builders are removed.
- **`New` is the only constructor**. `Default()` is the zero-config shortcut. Removed: `NewWithLoggerConfig`, `NewFromEnvironment`, `NewFromEnvSimple`, `NewSimple`, `NewEasyJSON`, `NewEasyJSONWithLevel`, `NewEasyBuilder`, `NewJSONLogger`, `NewTextLogger`, `NewSlogJSONLogger`, `NewSlogTextLogger`, `NewSlogLogger`, `NewStandardLogger`, `NewUnifiedLogger`, `NewWithLevel`, `NewWithLevelString`, `NewWithHandler`, `NewFromYAMLFile`, `NewFromYAMLEnv`, `LoadFromYAML`, `LoadFromYAMLData`, `LoadFromYAMLString`.
- **Trace ID retrieval renamed**: `GetTraceID(ctx) (string, bool)` → `TraceIDFromContext(ctx) string` (returns `""` when absent). Same renames for request and correlation IDs.
- **Removed package-level surface**: all handler builders (`AsyncHandler`, `BufferedHandler`, `ConditionalHandler`, `MiddlewareHandler`, `MultiHandler`, `RotatingHandler`, `HandlerBuilder`, `HandlerCompositor`, `HandlerRegistry`), all middleware (`CallerMiddleware`, `ContextExtractorMiddleware`, `LevelFilterMiddleware`, `LoggingMiddleware`, `MetricsMiddleware`, `RedactionMiddleware`, `SamplingMiddleware`, `StaticFieldsMiddleware`, `TimestampMiddleware`), all formatters (`CommonLogFormatter`, `ConsoleFormatter`, `JSONFormatter`, `TextFormatter`), all output types (`FileOutput`, `BufferedOutput`, `MultiOutput`, `WriterOutput`, `RotatingFileOutput`, `AsyncOutput`), the async-worker primitives, the redaction package (`RedactorChain`, `RegexRedactor`, `RedactedURL`, `RedactAPIKeys`), YAML config (`LoadFromYAML*`, `SaveToYAML`), Wire DI providers (`ProvideConfig`, `ProvideLogger`, etc.), HTTP middleware (`LogHTTPRequest`, `LogHTTPResponse`, `RequestLogger`, `TracingMiddleware`), context extractors (`BoolContextExtractor`, `StringContextExtractor`, etc.), and the UUID generator interface (`UUIDGenerator`, `SetUUIDGenerator`, `GetUUIDGenerator`, `MockUUIDGenerator`).
- **Shorthand removed**: `T()`, `D()`, `I()`, `E()` no longer exist.
- **Trace and Critical levels removed**: the four `slog` standard levels are the canonical set.
- **Examples directory removed**: docstring-level examples now live as `Example*` functions in `pkg/logging/example_test.go` and surface in godoc.
- **Configs directory removed**: YAML config support is gone; populate `Config` from your application's configuration system.
- **Dependencies removed**: `google/wire`, `go.uber.org/mock`, OpenTelemetry packages, `google/uuid`, and others. The module now declares zero external dependencies beyond the Go standard library.

### Added

- `WithContext(ctx, Logger) context.Context` — stash a Logger in a context.
- `FromContext(ctx) Logger` — retrieve the stashed Logger; falls back to the registered default (or built-in fallback) when missing. Never returns nil.
- `WithTraceFromContext(ctx, Logger) Logger` — convenience that attaches the trace ID from ctx as a `trace_id` field.
- `Logger.With(args ...any) Logger` — slog-style field attachment.
- `Logger.Enabled(ctx, slog.Level) bool` — for cheap guards around expensive log site work.
- `Default() Logger` — zero-config convenience.
- `Config.Handler` — escape hatch for callers that want a custom `slog.Handler`.

### Migration

Most call sites translate mechanically:

| v1 | v2 |
|---|---|
| `NewSimple()` | `Default()` |
| `NewJSONLogger(InfoLevel)` | `New(Config{Format: "json"})` |
| `NewTextLogger(InfoLevel)` | `New(Config{Format: "text"})` |
| `NewWithLoggerConfig(cfg)` | `New(Config{...})` |
| `logger.WithField("k", v)` | `logger.With("k", v)` |
| `logger.WithFields(m)` | `logger.With(slog.Any("k", v), ...)` |
| `GetTraceID(ctx)` (string, bool) | `TraceIDFromContext(ctx)` (string) |
| `logger.Critical("msg")` | `logger.Error("msg")` |
| `logger.Trace("msg")` | `logger.Debug("msg")` |
| YAML config | Populate `Config` from your config layer |
| Wire providers | Inline a `New(Config{...})` call in your composition root |
| HTTP middleware | Use your own; pair `WithContext`/`FromContext` with a thin handler |

## [1.0.0] - 2025-01-11

### Added

- Initial release of `go-logging`.

[Unreleased]: https://github.com/jedi-knights/go-logging/compare/v2.0.0...HEAD
[2.0.0]: https://github.com/jedi-knights/go-logging/releases/tag/v2.0.0
[1.0.0]: https://github.com/jedi-knights/go-logging/releases/tag/v1.0.0

# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Versions eligible for security patches:

| Version | Supported          |
| ------- | ------------------ |
| 2.0.x   | :white_check_mark: |
| < 2.0   | :x:                |

## Reporting a Vulnerability

We take the security of go-logging seriously. If you believe you have found a security vulnerability, please report it as described below.

**Do not create a public GitHub issue for security vulnerabilities.**

Report security vulnerabilities by email:

- **Email**: [INSERT SECURITY EMAIL]
- **Subject**: `[SECURITY] <brief description>`

### What to include

1. Description of the vulnerability
2. Steps to reproduce
3. Potential impact
4. Suggested fix (if you have one)
5. Your contact information for follow-up

### What to expect

- **Acknowledgment** within 48 hours of receipt
- **Initial assessment** within 5 business days
- **Regular updates** on progress
- **Target fix release** within 30 days for critical issues
- **Credit** in the security advisory, with your permission

## Security guidance for callers

This library is a thin wrapper over `log/slog`. The security posture of your application's logs depends on what you log and which `slog.Handler` you use, not on this library.

### Never log secrets

Treat passwords, tokens, API keys, session IDs, and PII as out-of-bounds for log output. Apply the rule at the call site:

```go
log.Info("login attempt", "user", username)             // ✅
log.Info("login attempt", "user", username, "pw", pw)   // ❌
```

### Redact at the handler, not at the call site

This library does not ship a redaction layer. If you need redaction (for example, to scrub anything that matches a secret pattern from your structured log output), wrap a `slog.Handler` that performs redaction and pass it via `Config.Handler`:

```go
log := logging.New(logging.Config{
    Handler: redaction.NewHandler(slog.NewJSONHandler(os.Stdout, nil)),
})
```

Wrapping at the handler ensures every record — including those emitted by code you do not own — is scrubbed before it lands on the wire.

### Be careful with HTTP headers

When logging request headers, log only an allowlist of safe headers. Authorization, Cookie, Set-Cookie, and any custom `X-Auth-*` headers should be redacted or omitted entirely.

### Be careful with errors

Errors frequently embed full connection strings, SQL fragments, and other internals. Wrap errors with a sanitized message at the boundary between your business logic and the log call:

```go
log.Error("database connection failed", "host", cfg.Host)  // ✅ sanitized
log.Error("database connection failed", "err", err)        // ❌ may include credentials
```

### Trace IDs are not secrets

`NewTraceID` returns RFC 4122 v4 UUIDs from `crypto/rand`. They are safe to log and propagate through context, but they should not be relied on for security (they are not authentication tokens).

## Known security considerations

### Thread safety

The `Logger` interface is safe for concurrent use. `With(args ...any)` returns a new Logger and does not mutate the receiver.

### Memory safety

No known memory leaks. The library holds no buffers and does no async work — every log call returns synchronously once the underlying `slog.Handler` returns.

### Dependency security

This library has zero external dependencies beyond the Go standard library. There is no transitive supply chain.

## Security updates

Security updates are released as:

- **Patch** versions for minor issues
- **Minor** versions for moderate issues that require additive API
- **Major** versions for breaking security changes

Subscribe to releases by watching the repository on GitHub.

## Disclosure timeline

We follow responsible disclosure:

1. **Day 0**: Vulnerability reported
2. **Day 2**: Report acknowledged
3. **Day 7**: Initial assessment and response plan
4. **Day 30**: Target fix release for critical issues
5. **Day 90**: Public disclosure, coordinated with reporter

## Questions

For questions about this security policy, contact:

- **Email**: [INSERT CONTACT EMAIL]
- **GitHub Discussions**: For general security questions

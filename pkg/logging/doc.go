// Package logging is a small slog-backed structured logging library.
//
// The package exposes a single Logger interface, a single Config struct, and
// a single constructor (New). Trace, request, and correlation IDs propagate
// through context, and a logger can be stashed in a context for retrieval by
// downstream code.
//
// # Quick start
//
// Construct a logger with explicit configuration and inject it:
//
//	l := logging.New(logging.Config{
//	    Level:       "info",
//	    Format:      "json",
//	    ServiceName: "auth-server",
//	    Environment: "production",
//	})
//	l.Info("server starting", "port", 8080)
//
// # Context-stored logger
//
// HTTP middleware typically constructs a request-scoped logger (enriched with
// the request's trace ID) and stashes it in the request context. Handlers
// retrieve it with FromContext:
//
//	func Handle(w http.ResponseWriter, r *http.Request) {
//	    log := logging.FromContext(r.Context())
//	    log.Info("handling request")
//	}
//
// FromContext never returns nil. When no logger has been stashed and no
// default has been registered via SetDefaultLogger, FromContext falls back
// to a package-internal info-level text logger writing to os.Stdout.
//
// # Trace IDs
//
// Trace, request, and correlation IDs are stored as context values and
// retrieved with the matching FromContext helpers:
//
//	ctx = logging.WithTraceID(ctx, logging.NewTraceID())
//	id := logging.TraceIDFromContext(ctx)  // returns "" when absent
//
// # Default logger
//
// The package exposes Debug, Info, Warn, and Error functions that delegate to
// a package-default Logger. This is a seam for tests and one-off scripts;
// production code should construct a Logger explicitly and inject it through
// the Logger interface rather than depend on the global. Install a default
// with SetDefaultLogger, or pass nil to detach the current default.
package logging

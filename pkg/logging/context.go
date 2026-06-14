package logging

import "context"

// contextKey is the private type used for all package-internal context keys.
// Using a named string type prevents collisions with other packages and keeps
// keys debug-friendly when contexts are printed.
type contextKey string

const (
	traceIDKey       contextKey = "logging.trace_id"
	requestIDKey     contextKey = "logging.request_id"
	correlationIDKey contextKey = "logging.correlation_id"
	loggerKey        contextKey = "logging.logger"
)

// WithContext returns a derived context carrying the given Logger.
// Retrieve the logger with FromContext.
func WithContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext returns the Logger stored in ctx, or the package default when no
// logger has been stashed. Never returns nil.
func FromContext(ctx context.Context) Logger {
	if ctx != nil {
		if l, ok := ctx.Value(loggerKey).(Logger); ok && l != nil {
			return l
		}
	}
	return defaultOrFallback()
}

// WithTraceFromContext returns a Logger enriched with the trace ID from ctx.
// When no trace ID is present, the original logger is returned unchanged.
func WithTraceFromContext(ctx context.Context, l Logger) Logger {
	id := TraceIDFromContext(ctx)
	if id == "" {
		return l
	}
	return l.With("trace_id", id)
}

package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// Config configures a Logger created via New.
//
// All fields are optional; the zero value is a valid info-level text logger
// writing to os.Stdout.
type Config struct {
	// Level is the minimum log level: "debug" | "info" | "warn" | "error".
	// Empty defaults to "info". Unknown values default to "info" and emit a
	// warning at construction time.
	Level string

	// Format is the output format: "json" | "text". Empty defaults to "text".
	Format string

	// Output is the destination writer. nil defaults to os.Stdout.
	Output io.Writer

	// ServiceName is attached to every record as the "service_name" attribute.
	ServiceName string

	// Environment is attached to every record as the "environment" attribute.
	Environment string

	// StaticFields are attached to every record.
	StaticFields map[string]any

	// Handler optionally overrides the slog.Handler. When set, Format and
	// Output are ignored — the provided handler is used directly. ServiceName,
	// Environment, and StaticFields still apply via With on the resulting logger.
	Handler slog.Handler
}

// New constructs a Logger from cfg.
func New(cfg Config) Logger {
	level, levelKnown := ParseLevel(cfg.Level)

	handler := cfg.Handler
	if handler == nil {
		out := cfg.Output
		if out == nil {
			out = os.Stdout
		}
		opts := &slog.HandlerOptions{Level: level}
		if strings.EqualFold(cfg.Format, "json") {
			handler = slog.NewJSONHandler(out, opts)
		} else {
			handler = slog.NewTextHandler(out, opts)
		}
	}

	inner := slog.New(handler)

	var attrs []any
	if cfg.ServiceName != "" {
		attrs = append(attrs, slog.String("service_name", cfg.ServiceName))
	}
	if cfg.Environment != "" {
		attrs = append(attrs, slog.String("environment", cfg.Environment))
	}
	for k, v := range cfg.StaticFields {
		attrs = append(attrs, slog.Any(k, v))
	}
	if len(attrs) > 0 {
		inner = inner.With(attrs...)
	}

	l := &slogLogger{inner: inner}

	if !levelKnown && cfg.Level != "" {
		l.Warn("unknown log level, defaulting to info", "level", cfg.Level)
	}

	return l
}

// Default is a convenience equivalent to New(Config{}).
func Default() Logger { return New(Config{}) }

// ParseLevel converts a level string to slog.Level. Comparison is case-insensitive.
// Returns (level, true) for known values ("debug", "info", "warn", "warning", "error", "").
// Returns (slog.LevelInfo, false) for unknown values.
func ParseLevel(s string) (slog.Level, bool) {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug, true
	case "info", "":
		return slog.LevelInfo, true
	case "warn", "warning":
		return slog.LevelWarn, true
	case "error":
		return slog.LevelError, true
	default:
		return slog.LevelInfo, false
	}
}

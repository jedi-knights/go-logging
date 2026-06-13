package logging_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/jedi-knights/go-logging/pkg/logging"
)

var uuidV4Pattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func TestWithTraceID_StoresAndRetrieves(t *testing.T) {
	t.Parallel()
	ctx := logging.WithTraceID(context.Background(), "trace-abc")
	if got := logging.TraceIDFromContext(ctx); got != "trace-abc" {
		t.Errorf("got %q, want %q", got, "trace-abc")
	}
}

func TestTraceIDFromContext_EmptyWhenAbsent(t *testing.T) {
	t.Parallel()
	if got := logging.TraceIDFromContext(context.Background()); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestWithRequestID_StoresAndRetrieves(t *testing.T) {
	t.Parallel()
	ctx := logging.WithRequestID(context.Background(), "req-xyz")
	if got := logging.RequestIDFromContext(ctx); got != "req-xyz" {
		t.Errorf("got %q, want %q", got, "req-xyz")
	}
}

func TestRequestIDFromContext_EmptyWhenAbsent(t *testing.T) {
	t.Parallel()
	if got := logging.RequestIDFromContext(context.Background()); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestWithCorrelationID_StoresAndRetrieves(t *testing.T) {
	t.Parallel()
	ctx := logging.WithCorrelationID(context.Background(), "corr-42")
	if got := logging.CorrelationIDFromContext(ctx); got != "corr-42" {
		t.Errorf("got %q, want %q", got, "corr-42")
	}
}

func TestCorrelationIDFromContext_EmptyWhenAbsent(t *testing.T) {
	t.Parallel()
	if got := logging.CorrelationIDFromContext(context.Background()); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestNewTraceID_MatchesUUIDv4Format(t *testing.T) {
	t.Parallel()
	id := logging.NewTraceID()
	if id == "" {
		t.Fatal("NewTraceID returned empty string")
	}
	if !uuidV4Pattern.MatchString(id) {
		t.Errorf("NewTraceID %q does not match UUID v4 format", id)
	}
}

func TestNewTraceID_UniquePerCall(t *testing.T) {
	t.Parallel()
	a := logging.NewTraceID()
	b := logging.NewTraceID()
	if a == b {
		t.Errorf("expected different IDs across calls, both = %q", a)
	}
}

func TestNewContextWithTrace_HasNonEmptyTraceID(t *testing.T) {
	t.Parallel()
	ctx := logging.NewContextWithTrace()
	if logging.TraceIDFromContext(ctx) == "" {
		t.Error("expected NewContextWithTrace to set a trace ID")
	}
}

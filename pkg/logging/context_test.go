package logging_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/jedi-knights/go-logging/pkg/logging"
)

func TestWithContext_StoresAndRetrievesLogger(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	original := logging.New(logging.Config{Output: &buf, Format: "text"})

	ctx := logging.WithContext(context.Background(), original)
	retrieved := logging.FromContext(ctx)

	retrieved.Info("hello")
	if !strings.Contains(buf.String(), "hello") {
		t.Errorf("expected retrieved logger to write to original buffer, got: %q", buf.String())
	}
}

func TestFromContext_ReturnsFallbackWhenLoggerMissing(t *testing.T) {
	t.Parallel()

	l := logging.FromContext(context.Background())
	if l == nil {
		t.Fatal("FromContext must never return nil; expected a fallback logger")
	}
	// Smoke check: calling a method should not panic.
	l.Info("fallback works")
}

func TestWithTraceFromContext_AttachesTraceIDFieldWhenPresent(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	base := logging.New(logging.Config{Output: &buf, Format: "text"})
	ctx := logging.WithTraceID(context.Background(), "trace-zzz")

	enriched := logging.WithTraceFromContext(ctx, base)
	enriched.Info("event")

	if !strings.Contains(buf.String(), "trace_id=trace-zzz") {
		t.Errorf("expected trace_id attribute in output, got: %q", buf.String())
	}
}

func TestWithTraceFromContext_NoOpWhenTraceMissing(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	base := logging.New(logging.Config{Output: &buf, Format: "text"})

	enriched := logging.WithTraceFromContext(context.Background(), base)
	enriched.Info("event")

	if strings.Contains(buf.String(), "trace_id=") {
		t.Errorf("expected no trace_id attribute when missing from context, got: %q", buf.String())
	}
}

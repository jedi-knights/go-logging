package logging_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/jedi-knights/go-logging/pkg/logging"
)

func TestNew_WritesInfoLineToConfiguredOutput(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logging.New(logging.Config{
		Level:  "info",
		Format: "text",
		Output: &buf,
	})
	l.Info("hello", "k", "v")

	out := buf.String()
	if !strings.Contains(out, "hello") {
		t.Fatalf("expected output to contain message %q, got: %q", "hello", out)
	}
	if !strings.Contains(out, "k=v") {
		t.Fatalf("expected output to contain attribute %q, got: %q", "k=v", out)
	}
}

func TestNew_JSONFormatEmitsValidJSON(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logging.New(logging.Config{
		Level:  "info",
		Format: "json",
		Output: &buf,
	})
	l.Info("hello", "k", "v")

	var record map[string]any
	if err := json.Unmarshal(buf.Bytes(), &record); err != nil {
		t.Fatalf("expected valid JSON, got error %v on output: %q", err, buf.String())
	}
	if record["msg"] != "hello" {
		t.Fatalf("expected msg=hello, got record: %v", record)
	}
	if record["k"] != "v" {
		t.Fatalf("expected k=v attribute, got record: %v", record)
	}
}

func TestNew_RespectsLevelFilter(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logging.New(logging.Config{
		Level:  "warn",
		Format: "text",
		Output: &buf,
	})
	l.Debug("debug-msg")
	l.Info("info-msg")
	l.Warn("warn-msg")
	l.Error("error-msg")

	out := buf.String()
	if strings.Contains(out, "debug-msg") {
		t.Errorf("debug should be filtered at warn level, got: %q", out)
	}
	if strings.Contains(out, "info-msg") {
		t.Errorf("info should be filtered at warn level, got: %q", out)
	}
	if !strings.Contains(out, "warn-msg") {
		t.Errorf("warn should be emitted at warn level, got: %q", out)
	}
	if !strings.Contains(out, "error-msg") {
		t.Errorf("error should be emitted at warn level, got: %q", out)
	}
}

func TestNew_WithReturnsLoggerCarryingFields(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	base := logging.New(logging.Config{
		Level:  "info",
		Format: "text",
		Output: &buf,
	})
	enriched := base.With("trace_id", "abc123")
	enriched.Info("event")

	out := buf.String()
	if !strings.Contains(out, "trace_id=abc123") {
		t.Fatalf("expected With-attached field in output, got: %q", out)
	}
}

func TestNew_ServiceNameAndEnvironmentAttachedToEveryRecord(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logging.New(logging.Config{
		Level:       "info",
		Format:      "json",
		Output:      &buf,
		ServiceName: "auth-server",
		Environment: "production",
	})
	l.Info("hello")

	var record map[string]any
	if err := json.Unmarshal(buf.Bytes(), &record); err != nil {
		t.Fatalf("expected valid JSON: %v", err)
	}
	if record["service_name"] != "auth-server" {
		t.Errorf("expected service_name attribute, got record: %v", record)
	}
	if record["environment"] != "production" {
		t.Errorf("expected environment attribute, got record: %v", record)
	}
}

func TestNew_EnabledReflectsConfiguredLevel(t *testing.T) {
	t.Parallel()

	l := logging.New(logging.Config{Level: "info"})
	ctx := context.Background()

	if l.Enabled(ctx, slog.LevelDebug) {
		t.Errorf("debug should not be enabled at info level")
	}
	if !l.Enabled(ctx, slog.LevelInfo) {
		t.Errorf("info should be enabled at info level")
	}
	if !l.Enabled(ctx, slog.LevelError) {
		t.Errorf("error should be enabled at info level")
	}
}

func TestNew_InfoContextEmitsLogLine(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logging.New(logging.Config{
		Level:  "info",
		Format: "text",
		Output: &buf,
	})
	l.InfoContext(context.Background(), "ctx-msg")

	if !strings.Contains(buf.String(), "ctx-msg") {
		t.Fatalf("expected InfoContext to emit message, got: %q", buf.String())
	}
}

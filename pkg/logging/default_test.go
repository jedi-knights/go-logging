package logging_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/jedi-knights/go-logging/pkg/logging"
)

// These tests mutate package-level state via SetDefaultLogger so they must NOT
// run with t.Parallel.

func TestSetDefaultLogger_PackageLevelInfoUsesIt(t *testing.T) {
	t.Cleanup(func() { logging.SetDefaultLogger(nil) })

	var buf bytes.Buffer
	l := logging.New(logging.Config{Output: &buf, Format: "text"})
	logging.SetDefaultLogger(l)

	logging.Info("hello-info")
	if !strings.Contains(buf.String(), "hello-info") {
		t.Errorf("expected package-level Info to delegate to SetDefaultLogger, got: %q", buf.String())
	}
}

func TestSetDefaultLogger_PackageLevelDebugWarnErrorUseIt(t *testing.T) {
	t.Cleanup(func() { logging.SetDefaultLogger(nil) })

	var buf bytes.Buffer
	l := logging.New(logging.Config{Level: "debug", Output: &buf, Format: "text"})
	logging.SetDefaultLogger(l)

	logging.Debug("dbg")
	logging.Warn("wrn")
	logging.Error("err")

	out := buf.String()
	for _, want := range []string{"dbg", "wrn", "err"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got: %q", want, out)
		}
	}
}

func TestSetDefaultLogger_NilResetsToBuiltInFallback(t *testing.T) {
	t.Cleanup(func() { logging.SetDefaultLogger(nil) })

	var buf bytes.Buffer
	logging.SetDefaultLogger(logging.New(logging.Config{Output: &buf, Format: "text"}))
	logging.SetDefaultLogger(nil)

	// After reset, package-level Info should not write to the previous buffer.
	logging.Info("post-reset")
	if strings.Contains(buf.String(), "post-reset") {
		t.Errorf("expected SetDefaultLogger(nil) to detach the previous default, got: %q", buf.String())
	}
}

func TestFromContext_FallsBackToDefaultLogger(t *testing.T) {
	t.Cleanup(func() { logging.SetDefaultLogger(nil) })

	var buf bytes.Buffer
	logging.SetDefaultLogger(logging.New(logging.Config{Output: &buf, Format: "text"}))

	l := logging.FromContext(context.Background())
	l.Info("fallback-routed")
	if !strings.Contains(buf.String(), "fallback-routed") {
		t.Errorf("expected FromContext fallback to route through SetDefaultLogger, got: %q", buf.String())
	}
}

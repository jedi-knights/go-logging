package logging_test

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/jedi-knights/go-logging/pkg/logging"
)

func ExampleNew() {
	var buf bytes.Buffer
	l := logging.New(logging.Config{
		Level:       "info",
		Format:      "text",
		Output:      &buf,
		ServiceName: "demo",
		Environment: "dev",
	})
	l.Info("hello", "user", "ada")

	// Strip the variable timestamp for a stable Output check.
	out := buf.String()
	if idx := strings.Index(out, "level="); idx >= 0 {
		out = out[idx:]
	}
	fmt.Print(out)
	// Output: level=INFO msg=hello service_name=demo environment=dev user=ada
}

func ExampleLogger_With() {
	var buf bytes.Buffer
	base := logging.New(logging.Config{Output: &buf, Format: "text"})

	requestLog := base.With("request_id", "r-42")
	requestLog.Info("starting")
	requestLog.Info("finished")

	for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		if idx := strings.Index(line, "level="); idx >= 0 {
			fmt.Println(line[idx:])
		}
	}
	// Output:
	// level=INFO msg=starting request_id=r-42
	// level=INFO msg=finished request_id=r-42
}

func ExampleWithContext() {
	var buf bytes.Buffer
	l := logging.New(logging.Config{Output: &buf, Format: "text"})

	ctx := logging.WithContext(context.Background(), l)
	logging.FromContext(ctx).Info("from-context")

	out := buf.String()
	if idx := strings.Index(out, "level="); idx >= 0 {
		out = out[idx:]
	}
	fmt.Print(out)
	// Output: level=INFO msg=from-context
}

func ExampleWithTraceFromContext() {
	var buf bytes.Buffer
	l := logging.New(logging.Config{Output: &buf, Format: "text"})

	ctx := logging.WithTraceID(context.Background(), "trace-42")
	logging.WithTraceFromContext(ctx, l).Info("event")

	out := buf.String()
	if idx := strings.Index(out, "level="); idx >= 0 {
		out = out[idx:]
	}
	fmt.Print(out)
	// Output: level=INFO msg=event trace_id=trace-42
}

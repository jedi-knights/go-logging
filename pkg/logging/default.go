package logging

import "sync"

// The package default logger is consulted by package-level Debug/Info/Warn/Error
// and by FromContext when no logger has been stashed in the context.
//
// SetDefaultLogger is intentionally narrow: it is a seam for tests and one-off
// scripts. Production code should construct a Logger explicitly via New and
// inject it through the Logger interface rather than reading from the global.
var (
	defaultMu     sync.RWMutex
	defaultLogger Logger
)

// fallbackLogger is the implicit default used when SetDefaultLogger has not been
// called (or has been called with nil). It is constructed once at package init
// so callers do not pay for repeated logger construction on the hot path.
var fallbackLogger = New(Config{})

// SetDefaultLogger installs l as the package default. Pass nil to detach the
// current default — subsequent package-level log calls then route to a built-in
// info-level text logger writing to os.Stdout.
func SetDefaultLogger(l Logger) {
	defaultMu.Lock()
	defaultLogger = l
	defaultMu.Unlock()
}

// defaultOrFallback returns the registered default, or the built-in fallback
// when none has been set.
func defaultOrFallback() Logger {
	defaultMu.RLock()
	l := defaultLogger
	defaultMu.RUnlock()
	if l != nil {
		return l
	}
	return fallbackLogger
}

// Debug logs at debug level via the package default.
func Debug(msg string, args ...any) { defaultOrFallback().Debug(msg, args...) }

// Info logs at info level via the package default.
func Info(msg string, args ...any) { defaultOrFallback().Info(msg, args...) }

// Warn logs at warn level via the package default.
func Warn(msg string, args ...any) { defaultOrFallback().Warn(msg, args...) }

// Error logs at error level via the package default.
func Error(msg string, args ...any) { defaultOrFallback().Error(msg, args...) }

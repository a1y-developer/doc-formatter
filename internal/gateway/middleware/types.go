package middleware

type contextKey struct {
	name string // name is the identifier for the context value.
}

var (
	APILoggerKey       = &contextKey{"api-logger"}
	RunLoggerKey       = &contextKey{"run-logger"}
	RunLoggerBufferKey = &contextKey{"run-logger-buffer"}
)

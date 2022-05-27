package routine

import "github.com/timandy/routine/bytesconv"

// RuntimeError an error type contains stack info.
type RuntimeError interface {
	// Message data when panic is raised.
	Message() Any

	// StackTrace stack when this instance is created.
	StackTrace() string

	// Error contains Message and StackTrace.
	Error() string
}

// NewRuntimeError create a new instance.
func NewRuntimeError(message Any) RuntimeError {
	return &runtimeError{message: message, stackTrace: bytesconv.String(traceStack())}
}

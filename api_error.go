package routine

// RuntimeError runtime error with stack info.
type RuntimeError interface {
	// Goid returns the goid of the coroutine that created the current error.
	Goid() int64

	// Gopc returns the pc of go statement that created the current error coroutine.
	Gopc() uintptr

	// Message returns the detail message string of this error.
	Message() string

	// StackTrace returns an array of stack trace elements, each representing one stack frame.
	StackTrace() []uintptr

	// Cause returns the cause of this error or nil if the cause is nonexistent or unknown.
	Cause() RuntimeError

	// Error returns a short description of this error.
	Error() string
}

// NewRuntimeError create a new RuntimeError instance.
func NewRuntimeError(cause any) RuntimeError {
	goid, gopc, msg, stackTrace, innerErr := runtimeErrorNew(cause)
	return &runtimeError{goid: goid, gopc: gopc, message: msg, stackTrace: stackTrace, cause: innerErr}
}

// NewRuntimeErrorWithMessage create a new RuntimeError instance.
func NewRuntimeErrorWithMessage(message string) RuntimeError {
	goid, gopc, msg, stackTrace, innerErr := runtimeErrorNewWithMessage(message)
	return &runtimeError{goid: goid, gopc: gopc, message: msg, stackTrace: stackTrace, cause: innerErr}
}

// NewRuntimeErrorWithMessageCause create a new RuntimeError instance.
func NewRuntimeErrorWithMessageCause(message string, cause any) RuntimeError {
	goid, gopc, msg, stackTrace, innerErr := runtimeErrorNewWithMessageCause(message, cause)
	return &runtimeError{goid: goid, gopc: gopc, message: msg, stackTrace: stackTrace, cause: innerErr}
}

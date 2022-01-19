package routine

type StackError interface {
	Message() Any
	StackTrace() string
	Error() string
}

func NewStackError(message Any) StackError {
	return &stackError{message: message, stackTrace: string(traceStack())}
}

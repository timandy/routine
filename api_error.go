package routine

import "fmt"

type StackError struct {
	error      Any
	stackTrace string
}

func NewStackError(error Any) *StackError {
	return &StackError{error: error, stackTrace: string(traceStack())}
}

func (fe *StackError) Error() string {
	return fmt.Sprintf("%v\n%v", fe.error, fe.stackTrace)
}

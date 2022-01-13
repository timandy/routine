package routine

import "fmt"

type StackError struct {
	error      Any
	stackTrace []byte
}

func stackError(error Any) *StackError {
	return &StackError{error: error, stackTrace: readStackBuf()}
}

func (fe *StackError) Error() string {
	return fmt.Sprintf("%v\n%v", fe.error, fe.stackTrace)
}

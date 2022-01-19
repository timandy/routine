package routine

import "fmt"

type stackError struct {
	message    Any
	stackTrace string
}

func (fe *stackError) Message() Any {
	return fe.message
}

func (fe *stackError) StackTrace() string {
	return fe.stackTrace
}

func (fe *stackError) Error() string {
	return fmt.Sprintf("%v\n%v", fe.message, fe.stackTrace)
}

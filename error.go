package routine

import "fmt"

type stackError struct {
	message    Any
	stackTrace string
}

func (se *stackError) Message() Any {
	return se.message
}

func (se *stackError) StackTrace() string {
	return se.stackTrace
}

func (se *stackError) Error() string {
	return fmt.Sprintf("%v\n%v", se.message, se.stackTrace)
}

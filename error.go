package routine

import "fmt"

type runtimeError struct {
	message    Any
	stackTrace string
}

func (re *runtimeError) Message() Any {
	return re.message
}

func (re *runtimeError) StackTrace() string {
	return re.stackTrace
}

func (re *runtimeError) Error() string {
	s := "RuntimeError"
	if message := fmt.Sprint(re.message); len(message) > 0 {
		s = s + ": " + message
	}
	return s + "\n" + re.stackTrace
}

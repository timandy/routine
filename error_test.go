package routine

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuntimeError_Goid(t *testing.T) {
	goid := Goid()
	err := NewRuntimeError(nil)
	assert.Equal(t, goid, err.Goid())
	fut := GoWait(func(token CancelToken) {
		assert.Equal(t, goid, err.Goid())
		assert.NotEqual(t, Goid(), err.Goid())
	})
	fut.Get()
}

func TestRuntimeError_Gopc(t *testing.T) {
	gopc := *getg().gopc
	err := NewRuntimeError(nil)
	assert.Equal(t, gopc, err.Gopc())
	fut := GoWait(func(token CancelToken) {
		assert.Equal(t, gopc, err.Gopc())
		assert.NotEqual(t, *getg().gopc, err.Gopc())
	})
	fut.Get()
}

func TestRuntimeError_Message(t *testing.T) {
	err := NewRuntimeError(nil)
	assert.Equal(t, "", err.Message())

	err2 := NewRuntimeError("Hello")
	assert.Equal(t, "Hello", err2.Message())

	err3 := NewRuntimeError(&person{Id: 1, Name: "Tim"})
	assert.Equal(t, "&{1 Tim}", err3.Message())
}

func TestRuntimeError_StackTrace(t *testing.T) {
	err := NewRuntimeError(nil)
	stackTrace := err.StackTrace()
	capturedStackTrace := captureStackTrace(0, 200)
	for i := 1; i < len(stackTrace); i++ {
		assert.Equal(t, capturedStackTrace[i], stackTrace[i])
	}
}

func TestRuntimeError_Cause(t *testing.T) {
	err := NewRuntimeError(nil)
	assert.Nil(t, err.Cause())

	err2 := NewRuntimeError(errors.New("error"))
	assert.Nil(t, err2.Cause())

	err3 := NewRuntimeError(&person{Id: 1, Name: "Tim"})
	assert.Nil(t, err3.Cause())

	err4 := NewRuntimeError(err)
	assert.Same(t, err, err4.Cause())
}

func TestRuntimeError_Error_EmptyMessage_NilError(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("", nil)
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, 5, len(lines))
	//
	line := lines[0]
	assert.Equal(t, "RuntimeError", line)
	//
	line = lines[1]
	assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestRuntimeError_Error_EmptyMessage_NilError() in "))
	assert.True(t, strings.HasSuffix(line, "error_test.go:68"))
	//
	line = lines[2]
	assert.True(t, strings.HasPrefix(line, "   at testing.tRunner() in "))
	//
	line = lines[3]
	assert.Equal(t, "   --- End of error stack trace ---", line)
	//
	line = lines[4]
	assert.True(t, strings.HasPrefix(line, "   created by testing.(*T).Run() in "))
}

func TestRuntimeError_Error_EmptyMessage_NormalError(t *testing.T) {
	cause := NewRuntimeError("this is inner error")
	err := NewRuntimeErrorWithMessageCause("", cause)
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, 9, len(lines))
	//
	line := lines[0]
	assert.Equal(t, "RuntimeError", line)
	//
	line = lines[1]
	assert.Equal(t, " ---> RuntimeError: this is inner error", line)
	//
	line = lines[2]
	assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestRuntimeError_Error_EmptyMessage_NormalError() in "))
	assert.True(t, strings.HasSuffix(line, "error_test.go:90"))
	//
	line = lines[3]
	assert.True(t, strings.HasPrefix(line, "   at testing.tRunner() in "))
	//
	line = lines[4]
	assert.Equal(t, "   --- End of inner error stack trace ---", line)
	//
	line = lines[5]
	assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestRuntimeError_Error_EmptyMessage_NormalError() in "))
	assert.True(t, strings.HasSuffix(line, "error_test.go:91"))
	//
	line = lines[6]
	assert.True(t, strings.HasPrefix(line, "   at testing.tRunner() in "))
	//
	line = lines[7]
	assert.Equal(t, "   --- End of error stack trace ---", line)
	//
	line = lines[8]
	assert.True(t, strings.HasPrefix(line, "   created by testing.(*T).Run() in "))
}

func TestRuntimeError_Error_NormalMessage_NilError(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("this is error message", nil)
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, 5, len(lines))
	//
	line := lines[0]
	assert.Equal(t, "RuntimeError: this is error message", line)
	//
	line = lines[1]
	assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestRuntimeError_Error_NormalMessage_NilError() in "))
	assert.True(t, strings.HasSuffix(line, "error_test.go:126"))
	//
	line = lines[2]
	assert.True(t, strings.HasPrefix(line, "   at testing.tRunner() in "))
	//
	line = lines[3]
	assert.Equal(t, "   --- End of error stack trace ---", line)
	//
	line = lines[4]
	assert.True(t, strings.HasPrefix(line, "   created by testing.(*T).Run() in "))
}

func TestRuntimeError_Error_NormalMessage_NormalError(t *testing.T) {
	cause := NewRuntimeError("this is inner error")
	err := NewRuntimeErrorWithMessageCause("this is error message", cause)
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, 9, len(lines))
	//
	line := lines[0]
	assert.Equal(t, "RuntimeError: this is error message", line)
	//
	line = lines[1]
	assert.Equal(t, " ---> RuntimeError: this is inner error", line)
	//
	line = lines[2]
	assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestRuntimeError_Error_NormalMessage_NormalError() in "))
	assert.True(t, strings.HasSuffix(line, "error_test.go:148"))
	//
	line = lines[3]
	assert.True(t, strings.HasPrefix(line, "   at testing.tRunner() in "))
	//
	line = lines[4]
	assert.Equal(t, "   --- End of inner error stack trace ---", line)
	//
	line = lines[5]
	assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestRuntimeError_Error_NormalMessage_NormalError() in "))
	assert.True(t, strings.HasSuffix(line, "error_test.go:149"))
	//
	line = lines[6]
	assert.True(t, strings.HasPrefix(line, "   at testing.tRunner() in "))
	//
	line = lines[7]
	assert.Equal(t, "   --- End of error stack trace ---", line)
	//
	line = lines[8]
	assert.True(t, strings.HasPrefix(line, "   created by testing.(*T).Run() in "))
}

func TestRuntimeError_Error_NilStackTrace(t *testing.T) {
	cause := NewRuntimeError("this is inner error")
	cause.(*runtimeError).stackTrace = nil
	err := NewRuntimeErrorWithMessageCause("this is error message", cause)
	err.(*runtimeError).stackTrace = nil
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, 5, len(lines))
	//
	line := lines[0]
	assert.Equal(t, "RuntimeError: this is error message", line)
	//
	line = lines[1]
	assert.Equal(t, " ---> RuntimeError: this is inner error", line)
	//
	line = lines[2]
	assert.Equal(t, "   --- End of inner error stack trace ---", line)
	//
	line = lines[3]
	assert.Equal(t, "   --- End of error stack trace ---", line)
	//
	line = lines[4]
	assert.True(t, strings.HasPrefix(line, "   created by testing.(*T).Run() in "))
}

func TestRuntimeError_Error_MainGoid(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("this is error message", nil)
	err.(*runtimeError).goid = 1
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, 3, len(lines))
	//
	line := lines[0]
	assert.Equal(t, "RuntimeError: this is error message", line)
	//
	line = lines[1]
	assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestRuntimeError_Error_MainGoid() in "))
	assert.True(t, strings.HasSuffix(line, "error_test.go:208"))
	//
	line = lines[2]
	assert.True(t, strings.HasPrefix(line, "   at testing.tRunner() in "))
}

func TestRuntimeError_Error_ZeroGopc(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("this is error message", nil)
	err.(*runtimeError).gopc = 0
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, 3, len(lines))
	//
	line := lines[0]
	assert.Equal(t, "RuntimeError: this is error message", line)
	//
	line = lines[1]
	assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestRuntimeError_Error_ZeroGopc() in "))
	assert.True(t, strings.HasSuffix(line, "error_test.go:225"))
	//
	line = lines[2]
	assert.True(t, strings.HasPrefix(line, "   at testing.tRunner() in "))
}

func TestArgumentNilError_Goid(t *testing.T) {
	goid := Goid()
	err := NewArgumentNilError("number", nil)
	assert.Equal(t, goid, err.Goid())
	fut := GoWait(func(token CancelToken) {
		assert.Equal(t, goid, err.Goid())
		assert.NotEqual(t, Goid(), err.Goid())
	})
	fut.Get()
}

func TestArgumentNilError_Gopc(t *testing.T) {
	gopc := *getg().gopc
	err := NewArgumentNilError("number", nil)
	assert.Equal(t, gopc, err.Gopc())
	fut := GoWait(func(token CancelToken) {
		assert.Equal(t, gopc, err.Gopc())
		assert.NotEqual(t, *getg().gopc, err.Gopc())
	})
	fut.Get()
}

func TestArgumentNilError_Message(t *testing.T) {
	err := NewArgumentNilError("", nil)
	assert.Equal(t, "Value cannot be null.", err.Message())

	err2 := NewArgumentNilError("", "Hello")
	assert.Equal(t, "Hello", err2.Message())

	err3 := NewArgumentNilError("number", nil)
	assert.Equal(t, "Value cannot be null."+newLine+"Parameter name: number.", err3.Message())

	err4 := NewArgumentNilError("number", "Hello")
	assert.Equal(t, "Hello"+newLine+"Parameter name: number.", err4.Message())
}

func TestArgumentNilError_StackTrace(t *testing.T) {
	err := NewArgumentNilError("", nil)
	stackTrace := err.StackTrace()
	capturedStackTrace := captureStackTrace(0, 200)
	for i := 1; i < len(stackTrace); i++ {
		assert.Equal(t, capturedStackTrace[i], stackTrace[i])
	}
}

func TestArgumentNilError_Cause(t *testing.T) {
	err := NewArgumentNilError("", nil)
	assert.Nil(t, err.Cause())

	err2 := NewArgumentNilError("", errors.New("error"))
	assert.Nil(t, err2.Cause())

	err3 := NewArgumentNilError("", &person{Id: 1, Name: "Tim"})
	assert.Nil(t, err3.Cause())

	err4 := NewArgumentNilError("", err)
	assert.Same(t, err, err4.Cause())
}

func TestArgumentNilError_Error(t *testing.T) {
	cause := NewRuntimeError("this is inner error")
	err := NewArgumentNilError("number", cause)
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, 10, len(lines))
	//
	line := lines[0]
	assert.Equal(t, "ArgumentNilError: Value cannot be null.", line)
	//
	line = lines[1]
	assert.Equal(t, "Parameter name: number.", line)
	//
	line = lines[2]
	assert.Equal(t, " ---> RuntimeError: this is inner error", line)
	//
	line = lines[3]
	assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestArgumentNilError_Error() in "))
	assert.True(t, strings.HasSuffix(line, "error_test.go:301"))
	//
	line = lines[4]
	assert.True(t, strings.HasPrefix(line, "   at testing.tRunner() in "))
	//
	line = lines[5]
	assert.Equal(t, "   --- End of inner error stack trace ---", line)
	//
	line = lines[6]
	assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestArgumentNilError_Error() in "))
	assert.True(t, strings.HasSuffix(line, "error_test.go:302"))
	//
	line = lines[7]
	assert.True(t, strings.HasPrefix(line, "   at testing.tRunner() in "))
	//
	line = lines[8]
	assert.Equal(t, "   --- End of error stack trace ---", line)
	//
	line = lines[9]
	assert.True(t, strings.HasPrefix(line, "   created by testing.(*T).Run() in "))
}

func TestArgumentNilError_ParamName(t *testing.T) {
	err := NewArgumentNilError("", nil)
	assert.Equal(t, "", err.ParamName())

	err2 := NewArgumentNilError("number", nil)
	assert.Equal(t, "number", err2.ParamName())
}

type ArgumentNilError struct {
	goid       int64
	gopc       uintptr
	message    string
	stackTrace []uintptr
	cause      RuntimeError
	paramName  string
}

func (ae *ArgumentNilError) Goid() int64 {
	return ae.goid
}

func (ae *ArgumentNilError) Gopc() uintptr {
	return ae.gopc
}

func (ae *ArgumentNilError) Message() string {
	builder := &strings.Builder{}
	if len(ae.message) == 0 {
		builder.WriteString("Value cannot be null.")
	} else {
		builder.WriteString(ae.message)
	}
	if len(ae.paramName) != 0 {
		builder.WriteString(newLine)
		builder.WriteString("Parameter name: ")
		builder.WriteString(ae.paramName)
		builder.WriteString(".")
	}
	return builder.String()
}

func (ae *ArgumentNilError) StackTrace() []uintptr {
	return ae.stackTrace
}

func (ae *ArgumentNilError) Cause() RuntimeError {
	return ae.cause
}

func (ae *ArgumentNilError) Error() string {
	return runtimeErrorError(ae)
}

func (ae *ArgumentNilError) ParamName() string {
	return ae.paramName
}

func NewArgumentNilError(paramName string, cause any) *ArgumentNilError {
	goid, gopc, msg, stackTrace, innerErr := runtimeErrorNew(cause)
	return &ArgumentNilError{goid: goid, gopc: gopc, message: msg, paramName: paramName, stackTrace: stackTrace, cause: innerErr}
}

package routine

import (
	"errors"
	"runtime"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRuntimeError_Nil(t *testing.T) {
	err := NewRuntimeError(nil)
	assertGoidGopc(t, err)
	assert.Equal(t, "", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeError_EmptyString(t *testing.T) {
	err := NewRuntimeError("")
	assertGoidGopc(t, err)
	assert.Equal(t, "", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeError_WhiteSpaceString(t *testing.T) {
	err := NewRuntimeError("\t")
	assertGoidGopc(t, err)
	assert.Equal(t, "\t", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: \t", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeError_NormalString(t *testing.T) {
	err := NewRuntimeError("this is error message")
	assertGoidGopc(t, err)
	assert.Equal(t, "this is error message", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: this is error message", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeError_NilError(t *testing.T) {
	var cause error
	err := NewRuntimeError(cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeError_NormalError(t *testing.T) {
	err := NewRuntimeError(errors.New("this is error message"))
	assertGoidGopc(t, err)
	assert.Equal(t, "this is error message", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: this is error message", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeError_NilRuntimeError(t *testing.T) {
	var cause RuntimeError
	err := NewRuntimeError(cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeError_NormalRuntimeError(t *testing.T) {
	cause := NewRuntimeError("this is inner message")
	err := NewRuntimeError(cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Same(t, cause, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError", lines[0])
	assert.Equal(t, " ---> RuntimeError: this is inner message", lines[1])
}

func TestNewRuntimeError_NilValue(t *testing.T) {
	var cause *person
	err := NewRuntimeError(cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "<nil>", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: <nil>", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeError_NormalValue(t *testing.T) {
	cause := person{Id: 1, Name: "Tim"}
	err := NewRuntimeError(cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "{1 Tim}", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: {1 Tim}", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessage_EmptyString(t *testing.T) {
	err := NewRuntimeErrorWithMessage("")
	assertGoidGopc(t, err)
	assert.Equal(t, "", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessage_WhiteSpaceString(t *testing.T) {
	err := NewRuntimeErrorWithMessage("\t")
	assertGoidGopc(t, err)
	assert.Equal(t, "\t", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: \t", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessage_NormalString(t *testing.T) {
	err := NewRuntimeErrorWithMessage("this is error message")
	assertGoidGopc(t, err)
	assert.Equal(t, "this is error message", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: this is error message", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_EmptyString_Nil(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("", nil)
	assertGoidGopc(t, err)
	assert.Equal(t, "", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_EmptyString_EmptyString(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("", "")
	assertGoidGopc(t, err)
	assert.Equal(t, "", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_EmptyString_WhiteSpaceString(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("", "\t")
	assertGoidGopc(t, err)
	assert.Equal(t, "\t", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: \t", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_EmptyString_NormalString(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("", "this is error message")
	assertGoidGopc(t, err)
	assert.Equal(t, "this is error message", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: this is error message", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_EmptyString_NilError(t *testing.T) {
	var cause error
	err := NewRuntimeErrorWithMessageCause("", cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_EmptyString_NormalError(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("", errors.New("this is error message"))
	assertGoidGopc(t, err)
	assert.Equal(t, "this is error message", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: this is error message", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_EmptyString_NilRuntimeError(t *testing.T) {
	var cause RuntimeError
	err := NewRuntimeErrorWithMessageCause("", cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_EmptyString_NormalRuntimeError(t *testing.T) {
	cause := NewRuntimeError("this is inner message")
	err := NewRuntimeErrorWithMessageCause("", cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Same(t, cause, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError", lines[0])
	assert.Equal(t, " ---> RuntimeError: this is inner message", lines[1])
}

func TestNewRuntimeErrorWithMessageCause_EmptyString_NilValue(t *testing.T) {
	var cause *person
	err := NewRuntimeErrorWithMessageCause("", cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "<nil>", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: <nil>", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_EmptyString_NormalValue(t *testing.T) {
	cause := person{Id: 1, Name: "Tim"}
	err := NewRuntimeErrorWithMessageCause("", cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "{1 Tim}", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: {1 Tim}", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_WhiteSpaceString_Nil(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("\t", nil)
	assertGoidGopc(t, err)
	assert.Equal(t, "\t", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: \t", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_WhiteSpaceString_EmptyString(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("\t", "")
	assertGoidGopc(t, err)
	assert.Equal(t, "\t", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: \t", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_WhiteSpaceString_WhiteSpaceString(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("\t", "\t")
	assertGoidGopc(t, err)
	assert.Equal(t, "\t - \t", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: \t - \t", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_WhiteSpaceString_NormalString(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("\t", "this is error message")
	assertGoidGopc(t, err)
	assert.Equal(t, "\t - this is error message", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: \t - this is error message", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_WhiteSpaceString_NilError(t *testing.T) {
	var cause error
	err := NewRuntimeErrorWithMessageCause("\t", cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "\t", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: \t", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_WhiteSpaceString_NormalError(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("\t", errors.New("this is error message"))
	assertGoidGopc(t, err)
	assert.Equal(t, "\t - this is error message", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: \t - this is error message", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_WhiteSpaceString_NilRuntimeError(t *testing.T) {
	var cause RuntimeError
	err := NewRuntimeErrorWithMessageCause("\t", cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "\t", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: \t", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_WhiteSpaceString_NormalRuntimeError(t *testing.T) {
	cause := NewRuntimeError("this is inner message")
	err := NewRuntimeErrorWithMessageCause("\t", cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "\t", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Same(t, cause, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: \t", lines[0])
	assert.Equal(t, " ---> RuntimeError: this is inner message", lines[1])
}

func TestNewRuntimeErrorWithMessageCause_WhiteSpaceString_NilValue(t *testing.T) {
	var cause *person
	err := NewRuntimeErrorWithMessageCause("\t", cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "\t - <nil>", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: \t - <nil>", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_WhiteSpaceString_NormalValue(t *testing.T) {
	cause := person{Id: 1, Name: "Tim"}
	err := NewRuntimeErrorWithMessageCause("\t", cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "\t - {1 Tim}", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: \t - {1 Tim}", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_NormalString_Nil(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("this is error message", nil)
	assertGoidGopc(t, err)
	assert.Equal(t, "this is error message", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: this is error message", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_NormalString_EmptyString(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("this is error message", "")
	assertGoidGopc(t, err)
	assert.Equal(t, "this is error message", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: this is error message", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_NormalString_WhiteSpaceString(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("this is error message", "\t")
	assertGoidGopc(t, err)
	assert.Equal(t, "this is error message - \t", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: this is error message - \t", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_NormalString_NormalString(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("this is error message", "this is error message")
	assertGoidGopc(t, err)
	assert.Equal(t, "this is error message - this is error message", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: this is error message - this is error message", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_NormalString_NilError(t *testing.T) {
	var cause error
	err := NewRuntimeErrorWithMessageCause("this is error message", cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "this is error message", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: this is error message", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_NormalString_NormalError(t *testing.T) {
	err := NewRuntimeErrorWithMessageCause("this is error message", errors.New("this is error message2"))
	assertGoidGopc(t, err)
	assert.Equal(t, "this is error message - this is error message2", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: this is error message - this is error message2", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_NormalString_NilRuntimeError(t *testing.T) {
	var cause RuntimeError
	err := NewRuntimeErrorWithMessageCause("this is error message", cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "this is error message", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: this is error message", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_NormalString_NormalRuntimeError(t *testing.T) {
	cause := NewRuntimeError("this is inner message")
	err := NewRuntimeErrorWithMessageCause("this is error message", cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "this is error message", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Same(t, cause, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: this is error message", lines[0])
	assert.Equal(t, " ---> RuntimeError: this is inner message", lines[1])
}

func TestNewRuntimeErrorWithMessageCause_NormalString_NilValue(t *testing.T) {
	var cause *person
	err := NewRuntimeErrorWithMessageCause("this is error message", cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "this is error message - <nil>", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: this is error message - <nil>", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func TestNewRuntimeErrorWithMessageCause_NormalString_NormalValue(t *testing.T) {
	cause := person{Id: 1, Name: "Tim"}
	err := NewRuntimeErrorWithMessageCause("this is error message", cause)
	assertGoidGopc(t, err)
	assert.Equal(t, "this is error message - {1 Tim}", err.Message())
	assert.Greater(t, len(err.StackTrace()), 1)
	assert.Nil(t, err.Cause())
	lines := strings.Split(err.Error(), newLine)
	assert.Equal(t, "RuntimeError: this is error message - {1 Tim}", lines[0])
	assert.Equal(t, "   at ", lines[1][:6])
}

func assertGoidGopc(t *testing.T, err RuntimeError) {
	assert.Equal(t, Goid(), err.Goid())
	assert.NotNil(t, runtime.FuncForPC(err.Gopc()-1))
}

//===

// BenchmarkDebugStack-4                             239652                 5305 ns/op         1024 B/op          1 allocs/op
func BenchmarkDebugStack(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = debug.Stack()
	}
}

// BenchmarkRuntimeError-4                           300091                 4020 ns/op         2484 B/op         15 allocs/op
func BenchmarkRuntimeError(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewRuntimeError(nil).Error()
	}
}

// BenchmarkRuntimeErrorWithMessage-4                302037                 3820 ns/op         2476 B/op         14 allocs/op
func BenchmarkRuntimeErrorWithMessage(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewRuntimeErrorWithMessage("").Error()
	}
}

// BenchmarkRuntimeErrorWithMessageCause-4           326098                 3679 ns/op         2652 B/op         14 allocs/op
func BenchmarkRuntimeErrorWithMessageCause(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewRuntimeErrorWithMessageCause("", nil).Error()
	}
}

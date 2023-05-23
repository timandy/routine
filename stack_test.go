package routine

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptureStackTrace(t *testing.T) {
	stackTrace := captureStackTrace(0, 10)
	assert.Greater(t, len(stackTrace), 2)
	frame, _ := runtime.CallersFrames(stackTrace).Next()
	assert.Equal(t, "github.com/timandy/routine.TestCaptureStackTrace", frame.Function)
	assert.Equal(t, 11, frame.Line)
	//
	stackTrace2 := captureStackSkip(1)
	assert.Greater(t, len(stackTrace2), 2)
	frame2, _ := runtime.CallersFrames(stackTrace2).Next()
	assert.Equal(t, "github.com/timandy/routine.TestCaptureStackTrace", frame2.Function)
	assert.Equal(t, 17, frame2.Line)
}

func TestCaptureStackTrace_Deep(t *testing.T) {
	stackTrace := captureStackDeep(20)
	assert.Greater(t, len(stackTrace), 20)
	frames := runtime.CallersFrames(stackTrace)
	//
	frame, more := frames.Next()
	assert.True(t, more)
	assert.Equal(t, "github.com/timandy/routine.captureStackDeepRecursive", frame.Function)
	assert.Equal(t, 93, frame.Line)
	//
	frame2, more2 := frames.Next()
	assert.True(t, more2)
	assert.Equal(t, "github.com/timandy/routine.captureStackDeepRecursive", frame2.Function)
	assert.Equal(t, 91, frame2.Line)
}

func TestCaptureStackTrace_Overflow(t *testing.T) {
	stackTrace := captureStackDeep(200)
	assert.Equal(t, 100, len(stackTrace))
}

func TestShowFrame(t *testing.T) {
	assert.False(t, showFrame("make"))
	assert.True(t, showFrame("strings.equal"))
	assert.True(t, showFrame("strings.Equal"))
	assert.False(t, showFrame("runtime.hello"))
	assert.True(t, showFrame("runtime.Hello"))
}

func TestSkipFrame(t *testing.T) {
	assert.False(t, skipFrame("runtime.a", true))
	assert.False(t, skipFrame("runtime.gopanic", true))
	assert.False(t, skipFrame("runtime.a", false))
	assert.True(t, skipFrame("runtime.gopanic", false))
}

func TestIsExportedRuntime(t *testing.T) {
	assert.False(t, isExportedRuntime(""))
	assert.False(t, isExportedRuntime("runtime."))
	assert.False(t, isExportedRuntime("hello_world"))
	assert.False(t, isExportedRuntime("runtime._"))
	assert.False(t, isExportedRuntime("runtime.a"))
	assert.True(t, isExportedRuntime("runtime.Hello"))
	assert.True(t, isExportedRuntime("runtime.Panic"))
}

func TestIsPanicRuntime(t *testing.T) {
	assert.False(t, isPanicRuntime(""))
	assert.False(t, isPanicRuntime("runtime."))
	assert.False(t, isPanicRuntime("hello_world"))
	assert.False(t, isPanicRuntime("runtime.a"))
	assert.True(t, isPanicRuntime("runtime.goPanicIndex"))
	assert.True(t, isPanicRuntime("runtime.gopanic"))
	assert.True(t, isPanicRuntime("runtime.panicshift"))
}

func captureStackSkip(skip int) []uintptr {
	return captureStackTrace(skip, 100)
}

func captureStackDeep(deep int) []uintptr {
	return captureStackDeepRecursive(1, deep)
}

func captureStackDeepRecursive(cur int, deep int) []uintptr {
	if cur < deep {
		cur++
		return captureStackDeepRecursive(cur, deep)
	}
	return captureStackTrace(0, 100)
}

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
	assert.Equal(t, 58, frame.Line)
	//
	frame2, more2 := frames.Next()
	assert.True(t, more2)
	assert.Equal(t, "github.com/timandy/routine.captureStackDeepRecursive", frame2.Function)
	assert.Equal(t, 56, frame2.Line)
}

func TestCaptureStackTrace_Overflow(t *testing.T) {
	stackTrace := captureStackDeep(200)
	assert.Equal(t, 100, len(stackTrace))
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

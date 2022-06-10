package routine

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestFuture_Complete(t *testing.T) {
	fut := NewFuture()
	go func() {
		fut.Complete(1)
	}()
	assert.Equal(t, 1, fut.Get())
}

func TestFuture_CompleteError_Common(t *testing.T) {
	defer func() {
		if cause := recover(); cause != nil {
			err := cause.(RuntimeError)
			assert.NotNil(t, err)
			assert.Equal(t, "1", err.Message())
			lines := strings.Split(err.Error(), newLine)
			//
			line := lines[0]
			assert.Equal(t, "RuntimeError: 1", line)
			//
			line = lines[1]
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.(*future).CompleteError() in "))
			assert.True(t, strings.HasSuffix(line, "future.go:20"))
			//
			line = lines[2]
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestFuture_CompleteError_Common."))
			assert.True(t, strings.HasSuffix(line, "future_test.go:50"))
			//
			line = lines[3]
			assert.True(t, strings.HasPrefix(line, "   at runtime.gopanic() in "))
			//
			line = lines[4]
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestFuture_CompleteError_Common."))
			assert.True(t, strings.HasSuffix(line, "future_test.go:53"))
		}
	}()
	//
	fut := NewFuture()
	go func() {
		defer func() {
			if cause := recover(); cause != nil {
				fut.CompleteError(cause)
			}
		}()
		panic(1)
	}()
	fut.Get()
	assert.Fail(t, "should not be here")
}

func TestFuture_CompleteError_RuntimeError(t *testing.T) {
	defer func() {
		if cause := recover(); cause != nil {
			err := cause.(RuntimeError)
			assert.NotNil(t, err)
			assert.Equal(t, "1", err.Message())
			lines := strings.Split(err.Error(), newLine)
			//
			line := lines[0]
			assert.Equal(t, "RuntimeError: 1", line)
			//
			line = lines[1]
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestFuture_CompleteError_RuntimeError."))
			assert.True(t, strings.HasSuffix(line, "future_test.go:87"))
			//
			line = lines[2]
			assert.True(t, strings.HasPrefix(line, "   at runtime.gopanic() in "))
			//
			line = lines[3]
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestFuture_CompleteError_RuntimeError."))
			assert.True(t, strings.HasSuffix(line, "future_test.go:90"))
		}
	}()
	//
	fut := NewFuture()
	go func() {
		defer func() {
			if cause := recover(); cause != nil {
				fut.CompleteError(NewRuntimeError(cause))
			}
		}()
		panic(1)
	}()
	fut.Get()
	assert.Fail(t, "should not be here")
}

func TestFuture_Get(t *testing.T) {
	run := false
	fut := NewFuture()
	go func() {
		time.Sleep(500 * time.Millisecond)
		run = true
		fut.Complete(nil)
	}()
	assert.Nil(t, fut.Get())
	assert.True(t, run)
}

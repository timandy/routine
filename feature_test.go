package routine

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFuture_Complete(t *testing.T) {
	fea := NewFeature()
	go func() {
		fea.Complete(1)
	}()
	assert.Equal(t, 1, fea.Get())
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
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.(*feature).CompleteError() in "))
			assert.True(t, strings.HasSuffix(line, "feature.go:20"))
			//
			line = lines[2]
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestFuture_CompleteError_Common."))
			assert.True(t, strings.HasSuffix(line, "feature_test.go:51"))
			//
			line = lines[3]
			assert.True(t, strings.HasPrefix(line, "   at runtime.gopanic() in "))
			//
			line = lines[4]
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestFuture_CompleteError_Common."))
			assert.True(t, strings.HasSuffix(line, "feature_test.go:54"))
		}
	}()
	//
	fea := NewFeature()
	go func() {
		defer func() {
			if cause := recover(); cause != nil {
				fea.CompleteError(cause)
			}
		}()
		panic(1)
	}()
	fea.Get()
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
			assert.True(t, strings.HasSuffix(line, "feature_test.go:88"))
			//
			line = lines[2]
			assert.True(t, strings.HasPrefix(line, "   at runtime.gopanic() in "))
			//
			line = lines[3]
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestFuture_CompleteError_RuntimeError."))
			assert.True(t, strings.HasSuffix(line, "feature_test.go:91"))
		}
	}()
	//
	fea := NewFeature()
	go func() {
		defer func() {
			if cause := recover(); cause != nil {
				fea.CompleteError(NewRuntimeError(cause))
			}
		}()
		panic(1)
	}()
	fea.Get()
	assert.Fail(t, "should not be here")
}

func TestFuture_Get(t *testing.T) {
	run := false
	fea := NewFeature()
	go func() {
		time.Sleep(500 * time.Millisecond)
		run = true
		fea.Complete(nil)
	}()
	assert.Nil(t, fea.Get())
	assert.True(t, run)
}

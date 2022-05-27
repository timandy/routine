package routine

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMessage(t *testing.T) {
	re := NewRuntimeError("Hello")
	assert.Equal(t, "Hello", re.Message())

	p := &person{Id: 1, Name: "Tim"}
	re2 := NewRuntimeError(p)
	assert.Same(t, p, re2.Message())
}

func TestStackTrace(t *testing.T) {
	re := NewRuntimeError("Hello")
	assert.True(t, strings.HasPrefix(re.StackTrace(), "goroutine "))
	assert.False(t, strings.HasSuffix(re.StackTrace(), "\n"))

	p := &person{Id: 1, Name: "Tim"}
	re2 := NewRuntimeError(p)
	assert.True(t, strings.HasPrefix(re2.StackTrace(), "goroutine "))
	assert.False(t, strings.HasSuffix(re2.StackTrace(), "\n"))
}

func TestError(t *testing.T) {
	err := NewRuntimeError(nil)
	assert.True(t, strings.HasPrefix(err.Error(), "RuntimeError: <nil>\ngoroutine "))
	assert.False(t, strings.HasSuffix(err.Error(), "\n"))
	//
	err2 := NewRuntimeError("")
	assert.True(t, strings.HasPrefix(err2.Error(), "RuntimeError\ngoroutine "))
	assert.False(t, strings.HasSuffix(err2.Error(), "\n"))
	//
	err3 := NewRuntimeError("Hello")
	assert.True(t, strings.HasPrefix(err3.Error(), "RuntimeError: Hello\ngoroutine "))
	assert.False(t, strings.HasSuffix(err3.Error(), "\n"))
	//
	defer func() {
		if msg := recover(); msg != nil {
			err4 := NewRuntimeError(msg)
			assert.True(t, strings.HasPrefix(err4.Error(), "RuntimeError: World\ngoroutine "))
			assert.False(t, strings.HasSuffix(err4.Error(), "\n"))
		}
	}()
	panic("World")
}

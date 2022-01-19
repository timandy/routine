package routine

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMessage(t *testing.T) {
	se := NewStackError("Hello")
	assert.Equal(t, "Hello", se.Message())

	p := &person{Id: 1, Name: "Tim"}
	se2 := NewStackError(p)
	assert.Same(t, p, se2.Message())
}

func TestStackTrace(t *testing.T) {
	se := NewStackError("Hello")
	assert.True(t, strings.HasPrefix(se.StackTrace(), "goroutine "))
	assert.False(t, strings.HasSuffix(se.StackTrace(), "\n"))

	p := &person{Id: 1, Name: "Tim"}
	se2 := NewStackError(p)
	assert.True(t, strings.HasPrefix(se2.StackTrace(), "goroutine "))
	assert.False(t, strings.HasSuffix(se2.StackTrace(), "\n"))
}

func TestError(t *testing.T) {
	err := NewStackError("Hello")
	assert.True(t, strings.HasPrefix(err.Error(), "Hello\ngoroutine "))
	assert.False(t, strings.HasSuffix(err.Error(), "\n"))
	//
	defer func() {
		if msg := recover(); msg != nil {
			err2 := NewStackError(msg)
			assert.True(t, strings.HasPrefix(err2.Error(), "World\ngoroutine "))
			assert.False(t, strings.HasSuffix(err2.Error(), "\n"))
		}
	}()
	panic("World")
}

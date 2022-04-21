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
	err := NewStackError(nil)
	assert.True(t, strings.HasPrefix(err.Error(), "StackError: <nil>\ngoroutine "))
	assert.False(t, strings.HasSuffix(err.Error(), "\n"))
	//
	err2 := NewStackError("")
	assert.True(t, strings.HasPrefix(err2.Error(), "StackError\ngoroutine "))
	assert.False(t, strings.HasSuffix(err2.Error(), "\n"))
	//
	err3 := NewStackError("Hello")
	assert.True(t, strings.HasPrefix(err3.Error(), "StackError: Hello\ngoroutine "))
	assert.False(t, strings.HasSuffix(err3.Error(), "\n"))
	//
	defer func() {
		if msg := recover(); msg != nil {
			err4 := NewStackError(msg)
			assert.True(t, strings.HasPrefix(err4.Error(), "StackError: World\ngoroutine "))
			assert.False(t, strings.HasSuffix(err4.Error(), "\n"))
		}
	}()
	panic("World")
}

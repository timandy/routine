package routine

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestStackError(t *testing.T) {
	err := NewStackError("Hello")
	assert.NotNil(t, err)
	//
	msg := err.Error()
	assert.True(t, strings.HasPrefix(msg, "Hello\ngoroutine "))
	assert.False(t, strings.HasSuffix(msg, "\n"))
	//
	defer func() {
		if e := recover(); e != nil {
			msg2 := NewStackError(e).Error()
			assert.True(t, strings.HasPrefix(msg2, "World\ngoroutine "))
			assert.False(t, strings.HasSuffix(msg2, "\n"))
		}
	}()
	panic("World")
}

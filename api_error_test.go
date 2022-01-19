package routine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStackError(t *testing.T) {
	se := NewStackError("Hello")
	assert.NotNil(t, se)
	//
	p, ok := se.(*stackError)
	assert.Same(t, p, se)
	assert.True(t, ok)
}

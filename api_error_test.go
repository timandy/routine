package routine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRuntimeError(t *testing.T) {
	re := NewRuntimeError("Hello")
	assert.NotNil(t, re)
	//
	p, ok := re.(*runtimeError)
	assert.Same(t, p, re)
	assert.True(t, ok)
}

package routine

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTraceStack(t *testing.T) {
	stack := traceStack()
	assert.Greater(t, len(stack), 64)
	assert.True(t, bytes.HasPrefix(stack, []byte("goroutine ")))
	assert.False(t, bytes.HasSuffix(stack, []byte("\n")))
}

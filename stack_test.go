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

func TestTraceStackDeep(t *testing.T) {
	stack := traceStackDeep(20)
	assert.Greater(t, len(stack), 1024)
	assert.True(t, bytes.HasPrefix(stack, []byte("goroutine ")))
	assert.False(t, bytes.HasSuffix(stack, []byte("\n")))
}

func traceStackDeep(deep int) []byte {
	return traceStackDeepRecursive(1, deep)
}

func traceStackDeepRecursive(cur int, tar int) []byte {
	if cur < tar {
		cur++
		return traceStackDeepRecursive(cur, tar)
	}
	return traceStack()
}

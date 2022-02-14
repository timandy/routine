package routine

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTraceTiny(t *testing.T) {
	tiny := traceTiny()
	assert.Equal(t, 64, len(tiny))
	assert.True(t, bytes.HasPrefix(tiny, anchor))
	assert.False(t, bytes.HasSuffix(tiny, []byte("\n")))
}

func TestTraceStack(t *testing.T) {
	stack := traceStack()
	assert.Greater(t, len(stack), 64)
	assert.True(t, bytes.HasPrefix(stack, anchor))
	assert.False(t, bytes.HasSuffix(stack, []byte("\n")))
}

func TestTraceAllStack(t *testing.T) {
	all := traceAllStack()
	assert.Greater(t, len(all), 64)
	assert.True(t, bytes.HasPrefix(all, anchor))
	assert.False(t, bytes.HasSuffix(all, []byte("\n")))
}

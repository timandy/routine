package routine

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTraceTiny(t *testing.T) {
	tiny := traceTiny()
	assert.Equal(t, 64, len(tiny))
	assert.True(t, strings.HasPrefix(string(tiny), "goroutine "))
	assert.False(t, strings.HasSuffix(string(tiny), "\n"))
}

func TestTraceStack(t *testing.T) {
	stack := traceStack()
	assert.Greater(t, len(stack), 64)
	assert.True(t, strings.HasPrefix(string(stack), "goroutine "))
	assert.False(t, strings.HasSuffix(string(stack), "\n"))
}

func TestTraceAllStack(t *testing.T) {
	all := traceAllStack()
	assert.Greater(t, len(all), 64)
	assert.True(t, strings.HasPrefix(string(all), "goroutine "))
	assert.False(t, strings.HasSuffix(string(all), "\n"))
}

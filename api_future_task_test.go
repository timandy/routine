package routine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCancelToken(t *testing.T) {
	task := NewFutureTask()
	token, ok := task.(CancelToken)
	assert.Same(t, task, token)
	assert.True(t, ok)
}

func TestNewFutureTask(t *testing.T) {
	task := NewFutureTask()
	assert.NotNil(t, task)
	//
	p, ok := task.(*futureTask)
	assert.Same(t, p, task)
	assert.True(t, ok)
}

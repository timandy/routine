package routine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCancelToken(t *testing.T) {
	task := NewFutureTask(func(task FutureTask) any { return nil })
	token, ok := task.(CancelToken)
	assert.Same(t, task, token)
	assert.True(t, ok)
}

func TestNewFutureTask(t *testing.T) {
	assert.Panics(t, func() {
		NewFutureTask(nil)
	})
	//
	task := NewFutureTask(func(task FutureTask) any { return nil })
	assert.NotNil(t, task)
	//
	p, ok := task.(*futureTask)
	assert.Same(t, p, task)
	assert.True(t, ok)
}

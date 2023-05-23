package routine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFutureCallable(t *testing.T) {
	var futureCallable FutureCallable = func(task FutureTask) interface{} {
		return "Hello"
	}
	assert.Equal(t, "Hello", futureCallable(nil))
	//
	var fun func(FutureTask) any = futureCallable
	assert.Equal(t, "Hello", fun(nil))
}

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

package routine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFutureCallable(t *testing.T) {
	var futureCallable FutureCallable[string] = func(task FutureTask[string]) string {
		return "Hello"
	}
	assert.Equal(t, "Hello", futureCallable(nil))
	//
	var fun func(FutureTask[string]) string = futureCallable
	assert.Equal(t, "Hello", fun(nil))
}

func TestCancelToken(t *testing.T) {
	task := NewFutureTask[any](func(task FutureTask[any]) any { return nil })
	token, ok := task.(CancelToken)
	assert.Same(t, task, token)
	assert.True(t, ok)
}

func TestNewFutureTask(t *testing.T) {
	assert.Panics(t, func() {
		NewFutureTask[any](nil)
	})
	//
	task := NewFutureTask[any](func(task FutureTask[any]) any { return nil })
	assert.NotNil(t, task)
	//
	p, ok := task.(*futureTask[any])
	assert.Same(t, p, task)
	assert.True(t, ok)
}

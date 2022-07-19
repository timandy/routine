package routine

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFuture_IsDone(t *testing.T) {
	fut := NewFuture()
	assert.False(t, fut.IsDone())
	fut.Complete(nil)
	assert.True(t, fut.IsDone())
	//
	fut2 := NewFuture()
	assert.False(t, fut2.IsDone())
	fut2.Cancel(nil)
	assert.True(t, fut2.IsDone())
	//
	fut3 := NewFuture()
	assert.False(t, fut3.IsDone())
	fut3.Fail(nil)
	assert.True(t, fut3.IsDone())
}

func TestFuture_IsCanceled(t *testing.T) {
	fut := NewFuture()
	assert.False(t, fut.IsCanceled())
	fut.Cancel(nil)
	assert.True(t, fut.IsCanceled())
}

func TestFuture_IsFailed(t *testing.T) {
	fut := NewFuture()
	assert.False(t, fut.IsFailed())
	fut.Fail(nil)
	assert.True(t, fut.IsFailed())
}

func TestFuture_Complete_AfterCancel(t *testing.T) {
	fut := NewFuture()
	go func() {
		fut.Cancel(1)
	}()
	assert.Panics(t, func() {
		fut.Get()
	})
	assert.True(t, fut.IsCanceled())
	//
	go func() {
		fut.Complete(2)
	}()
	assert.Panics(t, func() {
		fut.Get()
	})
	assert.True(t, fut.IsCanceled())
}

func TestFuture_Complete_Common(t *testing.T) {
	fut := NewFuture()
	go func() {
		fut.Complete(1)
	}()
	assert.Equal(t, 1, fut.Get())
	//complete again won't change the result
	go func() {
		fut.Complete(2)
	}()
	assert.Equal(t, 1, fut.Get())
}

func TestFuture_Cancel_AfterComplete(t *testing.T) {
	fut := NewFuture()
	go func() {
		fut.Complete(1)
	}()
	assert.Equal(t, 1, fut.Get())
	fut.Cancel(1)
	assert.False(t, fut.IsCanceled())
	assert.Equal(t, 1, fut.Get())
}

func TestFuture_Cancel_Common(t *testing.T) {
	fut := NewFuture()
	go func() {
		fut.Cancel("Canceled")
	}()
	assert.Panics(t, func() {
		fut.Get()
	})
	assert.True(t, fut.IsCanceled())
	assert.Equal(t, "Canceled", fut.(*future).error.Message())
	assert.Nil(t, fut.(*future).error.Cause())
}

func TestFuture_Cancel_RuntimeError(t *testing.T) {
	fut3 := NewFuture()
	go func() {
		fut3.Cancel(NewRuntimeError("Canceled2"))
	}()
	assert.Panics(t, func() {
		fut3.Get()
	})
	assert.True(t, fut3.IsCanceled())
	assert.Equal(t, "Canceled2", fut3.(*future).error.Message())
	assert.Nil(t, fut3.(*future).error.Cause())
}

func TestFuture_Fail_AfterComplete(t *testing.T) {
	fut := NewFuture()
	go func() {
		fut.Complete(1)
	}()
	assert.Equal(t, 1, fut.Get())
	fut.Fail(1)
	assert.False(t, fut.IsFailed())
	assert.Equal(t, 1, fut.Get())
}

func TestFuture_Fail_Common(t *testing.T) {
	defer func() {
		if cause := recover(); cause != nil {
			err := cause.(RuntimeError)
			assert.NotNil(t, err)
			assert.Equal(t, "1", err.Message())
			lines := strings.Split(err.Error(), newLine)
			//
			line := lines[0]
			assert.Equal(t, "RuntimeError: 1", line)
			//
			line = lines[1]
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.(*future).Fail() in "))
			assert.True(t, strings.HasSuffix(line, "future.go:78"))
			//
			line = lines[2]
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestFuture_Fail_Common."))
			assert.True(t, strings.HasSuffix(line, "future_test.go:154"))
			//
			line = lines[3]
			assert.True(t, strings.HasPrefix(line, "   at runtime.gopanic() in "))
			//
			line = lines[4]
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestFuture_Fail_Common."))
			assert.True(t, strings.HasSuffix(line, "future_test.go:157"))
		}
	}()
	//
	fut := NewFuture()
	go func() {
		defer func() {
			if cause := recover(); cause != nil {
				fut.Fail(cause)
			}
		}()
		panic(1)
	}()
	fut.Get()
	assert.Fail(t, "should not be here")
}

func TestFuture_Fail_RuntimeError(t *testing.T) {
	defer func() {
		if cause := recover(); cause != nil {
			err := cause.(RuntimeError)
			assert.NotNil(t, err)
			assert.Equal(t, "1", err.Message())
			lines := strings.Split(err.Error(), newLine)
			//
			line := lines[0]
			assert.Equal(t, "RuntimeError: 1", line)
			//
			line = lines[1]
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestFuture_Fail_RuntimeError."))
			assert.True(t, strings.HasSuffix(line, "future_test.go:191"))
			//
			line = lines[2]
			assert.True(t, strings.HasPrefix(line, "   at runtime.gopanic() in "))
			//
			line = lines[3]
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestFuture_Fail_RuntimeError."))
			assert.True(t, strings.HasSuffix(line, "future_test.go:194"))
		}
	}()
	//
	fut := NewFuture()
	go func() {
		defer func() {
			if cause := recover(); cause != nil {
				fut.Fail(NewRuntimeError(cause))
			}
		}()
		panic(1)
	}()
	fut.Get()
	assert.Fail(t, "should not be here")
}

func TestFuture_Get_Nil(t *testing.T) {
	run := false
	fut := NewFuture()
	go func() {
		time.Sleep(500 * time.Millisecond)
		run = true
		fut.Complete(nil)
	}()
	assert.Nil(t, fut.Get())
	assert.True(t, run)
}

func TestFuture_Get_Common(t *testing.T) {
	run := false
	fut := NewFuture()
	go func() {
		time.Sleep(500 * time.Millisecond)
		run = true
		fut.Complete(1)
	}()
	assert.Equal(t, 1, fut.Get())
	assert.True(t, run)
}

func TestFuture_GetWithTimeout_Complete(t *testing.T) {
	run := false
	fut := NewFuture()
	go func() {
		if fut.IsCanceled() {
			return
		}
		run = true
		fut.Complete(1)
	}()
	assert.Equal(t, 1, fut.GetWithTimeout(500*time.Millisecond))
	assert.True(t, run)
}

func TestFuture_GetWithTimeout_Fail(t *testing.T) {
	run := false
	fut := NewFuture()
	go func() {
		if fut.IsCanceled() {
			return
		}
		run = true
		fut.Fail(1)
	}()
	assert.Panics(t, func() {
		fut.GetWithTimeout(200 * time.Millisecond)
	})
	assert.True(t, run)
	//
	assert.True(t, fut.IsFailed())
	assert.Equal(t, "1", fut.(*future).error.Message())
	assert.Nil(t, fut.(*future).error.Cause())
}

func TestFuture_GetWithTimeout_Timeout(t *testing.T) {
	run := false
	fut := NewFuture()
	go func() {
		time.Sleep(500 * time.Millisecond)
		if fut.IsCanceled() {
			return
		}
		run = true
		fut.Complete(nil)
	}()
	assert.Panics(t, func() {
		fut.GetWithTimeout(1 * time.Millisecond)
	})
	assert.False(t, run)
	//
	assert.True(t, fut.IsCanceled())
	assert.Equal(t, "Task execution timeout after 1ms.", fut.(*future).error.Message())
	assert.Nil(t, fut.(*future).error.Cause())
}

func TestFuture_Routine_Complete(t *testing.T) {
	var fut Future
	fut = GoWaitResult(func() Any {
		for i := 0; i < 5; i++ {
			if fut.IsCanceled() {
				panic("canceled")
			}
			time.Sleep(10 * time.Millisecond)
		}
		return 1
	})
	assert.Equal(t, 1, fut.GetWithTimeout(100*time.Millisecond))
}

func TestFuture_Routine_Fail(t *testing.T) {
	var fut Future
	fut = GoWaitResult(func() Any {
		panic("something error")
	})
	assert.Panics(t, func() {
		fut.GetWithTimeout(10 * time.Millisecond)
	})
	assert.True(t, fut.IsFailed())
	assert.Equal(t, "something error", fut.(*future).error.Message())
	assert.Nil(t, fut.(*future).error.Cause())
}

func TestFuture_Routine_Timeout(t *testing.T) {
	var fut Future
	fut = GoWaitResult(func() Any {
		for i := 0; i < 5; i++ {
			if fut.IsCanceled() {
				panic("canceled")
			}
			time.Sleep(10 * time.Millisecond)
		}
		return 1
	})
	assert.Panics(t, func() {
		fut.GetWithTimeout(1 * time.Millisecond)
	})
	assert.True(t, fut.IsCanceled())
	assert.Equal(t, "Task execution timeout after 1ms.", fut.(*future).error.Message())
	assert.Nil(t, fut.(*future).error.Cause())
}

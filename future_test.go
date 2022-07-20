package routine

import (
	"strings"
	"sync"
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
	fut2.Cancel()
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
	fut.Cancel()
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
		fut.Cancel()
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
	fut.Cancel()
	assert.False(t, fut.IsCanceled())
	assert.Equal(t, 1, fut.Get())
}

func TestFuture_Cancel_Common(t *testing.T) {
	fut := NewFuture()
	go func() {
		fut.Cancel()
	}()
	assert.Panics(t, func() {
		fut.Get()
	})
	assert.True(t, fut.IsCanceled())
	assert.Equal(t, "Task was canceled.", fut.(*future).error.Message())
	assert.Nil(t, fut.(*future).error.Cause())
}

func TestFuture_Cancel_RuntimeError(t *testing.T) {
	fut3 := NewFuture()
	go func() {
		fut3.Cancel()
	}()
	assert.Panics(t, func() {
		fut3.Get()
	})
	assert.True(t, fut3.IsCanceled())
	assert.Equal(t, "Task was canceled.", fut3.(*future).error.Message())
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
			assert.True(t, strings.HasSuffix(line, "future.go:74"))
			//
			line = lines[2]
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestFuture_Fail_Common."))
			assert.True(t, strings.HasSuffix(line, "future_test.go:155"))
			//
			line = lines[3]
			assert.True(t, strings.HasPrefix(line, "   at runtime.gopanic() in "))
			//
			line = lines[4]
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestFuture_Fail_Common."))
			assert.True(t, strings.HasSuffix(line, "future_test.go:158"))
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
			assert.True(t, strings.HasSuffix(line, "future_test.go:192"))
			//
			line = lines[2]
			assert.True(t, strings.HasPrefix(line, "   at runtime.gopanic() in "))
			//
			line = lines[3]
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestFuture_Fail_RuntimeError."))
			assert.True(t, strings.HasSuffix(line, "future_test.go:195"))
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
		time.Sleep(100 * time.Millisecond)
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
		time.Sleep(100 * time.Millisecond)
		run = true
		fut.Complete(1)
	}()
	assert.Equal(t, 1, fut.Get())
	assert.True(t, run)
}

func TestFuture_GetWithTimeout_Complete(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	run := false
	fut := NewFuture()
	go func() {
		defer wg.Done()
		//
		if fut.IsCanceled() {
			return
		}
		run = true
		fut.Complete(1)
	}()
	assert.Equal(t, 1, fut.GetWithTimeout(100*time.Millisecond))
	assert.True(t, run)
	//
	wg.Wait()
}

func TestFuture_GetWithTimeout_Fail(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	run := false
	fut := NewFuture()
	go func() {
		defer wg.Done()
		//
		if fut.IsCanceled() {
			return
		}
		run = true
		fut.Fail(1)
	}()
	assert.Panics(t, func() {
		fut.GetWithTimeout(100 * time.Millisecond)
	})
	assert.True(t, run)
	//
	assert.True(t, fut.IsFailed())
	assert.Equal(t, "1", fut.(*future).error.Message())
	assert.Nil(t, fut.(*future).error.Cause())
	//
	wg.Wait()
}

func TestFuture_GetWithTimeout_Timeout(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	run := false
	fut := NewFuture()
	go func() {
		defer wg.Done()
		//
		time.Sleep(100 * time.Millisecond)
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
	//
	wg.Wait()
}

func TestFuture_Routine_Complete(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	fut := GoWaitResult(func(token CancelToken) any {
		defer wg.Done()
		//
		if token.IsCanceled() {
			panic("canceled")
		}
		time.Sleep(1 * time.Millisecond)
		return 1
	})
	assert.Equal(t, 1, fut.GetWithTimeout(100*time.Millisecond))
	//
	wg.Wait()
}

func TestFuture_Routine_Cancel(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	fut := GoWaitResult(func(token CancelToken) any {
		defer wg.Done()
		//
		token.Cancel()
		return 1
	})
	assert.Panics(t, func() {
		fut.GetWithTimeout(100 * time.Millisecond)
	})
	assert.True(t, fut.IsCanceled())
	assert.Equal(t, "Task was canceled.", fut.(*future).error.Message())
	assert.Nil(t, fut.(*future).error.Cause())
	//
	wg.Wait()
}

func TestFuture_Routine_CancelInParent(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	finished := false
	fut := GoWaitResult(func(token CancelToken) any {
		defer wg.Done()
		//
		for i := 0; i < 10; i++ {
			time.Sleep(10 * time.Millisecond)
			if token.IsCanceled() {
				return 0
			}
		}
		finished = true
		return 1
	})
	fut.Cancel()
	//
	wg.Wait()
	//
	assert.False(t, finished)
	assert.True(t, fut.IsCanceled())
	assert.Equal(t, "Task was canceled.", fut.(*future).error.Message())
	assert.Nil(t, fut.(*future).error.Cause())
}

func TestFuture_Routine_Fail(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	fut := GoWaitResult(func(token CancelToken) any {
		defer wg.Done()
		//
		if token.IsCanceled() {
			return 1
		}
		panic("something error")
	})
	assert.Panics(t, func() {
		fut.GetWithTimeout(100 * time.Millisecond)
	})
	assert.True(t, fut.IsFailed())
	assert.Equal(t, "something error", fut.(*future).error.Message())
	assert.Nil(t, fut.(*future).error.Cause())
	//
	wg.Wait()
}

func TestFuture_Routine_Timeout(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	fut := GoWaitResult(func(token CancelToken) any {
		defer wg.Done()
		//
		for i := 0; i < 10; i++ {
			if token.IsCanceled() {
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
	//
	wg.Wait()
}

func TestFuture_Routine_TimeoutThenComplete(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	fut := GoWait(func(token CancelToken) {
		defer wg.Done()
		//
		ft := token.(*future)
		ft.lock.Lock()
		defer ft.lock.Unlock()
		ft.result = 1
		ft.status = completed
		time.Sleep(50 * time.Millisecond)
		ft.await.Done()
	})
	assert.Equal(t, 1, fut.GetWithTimeout(10*time.Millisecond))
	assert.Equal(t, 1, fut.Get())
	//
	wg.Wait()
}

func TestFuture_Routine_TimeoutThenCancel(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	fut := GoWait(func(token CancelToken) {
		defer wg.Done()
		//
		ft := token.(*future)
		ft.lock.Lock()
		defer ft.lock.Unlock()
		ft.error = NewRuntimeError("canceled.")
		ft.status = canceled
		time.Sleep(50 * time.Millisecond)
		ft.await.Done()
	})
	assert.Panics(t, func() {
		fut.GetWithTimeout(10 * time.Millisecond)
	})
	//
	assert.True(t, fut.IsCanceled())
	assert.Equal(t, "canceled.", fut.(*future).error.Message())
	assert.Nil(t, fut.(*future).error.Cause())
	assert.Panics(t, func() {
		fut.Get()
	})
	//
	wg.Wait()
}

func TestFuture_Routine_TimeoutThenFail(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	fut := GoWait(func(token CancelToken) {
		defer wg.Done()
		//
		ft := token.(*future)
		ft.lock.Lock()
		defer ft.lock.Unlock()
		ft.error = NewRuntimeError("failed.")
		ft.status = failed
		time.Sleep(50 * time.Millisecond)
		ft.await.Done()
	})
	assert.Panics(t, func() {
		fut.GetWithTimeout(10 * time.Millisecond)
	})
	//
	assert.True(t, fut.IsFailed())
	assert.Equal(t, "failed.", fut.(*future).error.Message())
	assert.Nil(t, fut.(*future).error.Cause())
	assert.Panics(t, func() {
		fut.Get()
	})
	//
	wg.Wait()
}

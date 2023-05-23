package routine

import (
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFutureTask_IsDone(t *testing.T) {
	task := NewFutureTask(func(task FutureTask) any { return nil })
	assert.False(t, task.IsDone())
	task.Complete(nil)
	assert.True(t, task.IsDone())
	//
	task2 := NewFutureTask(func(task FutureTask) any { return nil })
	assert.False(t, task2.IsDone())
	task2.Cancel()
	assert.True(t, task2.IsDone())
	//
	task3 := NewFutureTask(func(task FutureTask) any { return nil })
	assert.False(t, task3.IsDone())
	task3.Fail(nil)
	assert.True(t, task3.IsDone())
	//
	task4 := NewFutureTask(func(task FutureTask) any { return nil })
	assert.False(t, task4.IsDone())
	task4.(*futureTask).state = taskStateRunning
	assert.False(t, task4.IsDone())
}

func TestFutureTask_IsCanceled(t *testing.T) {
	task := NewFutureTask(func(task FutureTask) any { return nil })
	assert.False(t, task.IsCanceled())
	task.Cancel()
	assert.True(t, task.IsCanceled())
}

func TestFutureTask_IsFailed(t *testing.T) {
	task := NewFutureTask(func(task FutureTask) any { return nil })
	assert.False(t, task.IsFailed())
	task.Fail(nil)
	assert.True(t, task.IsFailed())
}

func TestFutureTask_Complete_AfterCancel(t *testing.T) {
	task := NewFutureTask(func(task FutureTask) any { return nil })
	go func() {
		task.Cancel()
	}()
	assert.Panics(t, func() {
		task.Get()
	})
	assert.True(t, task.IsCanceled())
	//
	task.Complete(2)
	assert.Panics(t, func() {
		task.Get()
	})
	assert.True(t, task.IsCanceled())
}

func TestFutureTask_Complete_AfterComplete(t *testing.T) {
	task := NewFutureTask(func(task FutureTask) any { return 1 })
	task.Run()
	assert.Equal(t, 1, task.Get())
	task.Complete(2)
	assert.Equal(t, 1, task.Get())
	//
	run := false
	task2 := NewFutureTask(func(task FutureTask) any {
		run = true
		return 1
	})
	task2.Complete(2)
	task2.Run()
	assert.Equal(t, 2, task2.Get())
	assert.False(t, run)
}

func TestFutureTask_Complete_Common(t *testing.T) {
	task := NewFutureTask(func(task FutureTask) any { return nil })
	go func() {
		task.Complete(1)
	}()
	assert.Equal(t, 1, task.Get())
	//complete again won't change the result
	go func() {
		task.Complete(2)
	}()
	assert.Equal(t, 1, task.Get())
}

func TestFutureTask_Cancel_AfterComplete(t *testing.T) {
	task := NewFutureTask(func(task FutureTask) any { return nil })
	go func() {
		task.Complete(1)
	}()
	assert.Equal(t, 1, task.Get())
	task.Cancel()
	assert.False(t, task.IsCanceled())
	assert.Equal(t, 1, task.Get())
}

func TestFutureTask_Cancel_Common(t *testing.T) {
	task := NewFutureTask(func(task FutureTask) any { return nil })
	go func() {
		task.Cancel()
	}()
	assert.Panics(t, func() {
		task.Get()
	})
	assert.True(t, task.IsCanceled())
	assert.Equal(t, "Task was canceled.", task.(*futureTask).error.Message())
	assert.Nil(t, task.(*futureTask).error.Cause())
}

func TestFutureTask_Cancel_RuntimeError(t *testing.T) {
	task3 := NewFutureTask(func(task FutureTask) any { return nil })
	go func() {
		task3.Cancel()
	}()
	assert.Panics(t, func() {
		task3.Get()
	})
	assert.True(t, task3.IsCanceled())
	assert.Equal(t, "Task was canceled.", task3.(*futureTask).error.Message())
	assert.Nil(t, task3.(*futureTask).error.Cause())
}

func TestFutureTask_Fail_AfterComplete(t *testing.T) {
	task := NewFutureTask(func(task FutureTask) any { return nil })
	go func() {
		task.Complete(1)
	}()
	assert.Equal(t, 1, task.Get())
	task.Fail(1)
	assert.False(t, task.IsFailed())
	assert.Equal(t, 1, task.Get())
}

func TestFutureTask_Fail_Common(t *testing.T) {
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
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestFutureTask_Fail_Common."))
			assert.True(t, strings.HasSuffix(line, "future_task_test.go:169"))
		}
	}()
	//
	task := NewFutureTask(func(task FutureTask) any { return nil })
	go func() {
		defer func() {
			if cause := recover(); cause != nil {
				task.Fail(cause)
			}
		}()
		panic(1)
	}()
	task.Get()
	assert.Fail(t, "should not be here")
}

func TestFutureTask_Fail_RuntimeError(t *testing.T) {
	defer func() {
		if cause := recover(); cause != nil {
			err := cause.(RuntimeError)
			assert.NotNil(t, err)
			assert.Equal(t, "1", err.Message())
			lines := strings.Split(err.Error(), newLine)
			assert.Equal(t, 4, len(lines))
			//
			line := lines[0]
			assert.Equal(t, "RuntimeError: 1", line)
			//
			line = lines[1]
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestFutureTask_Fail_RuntimeError."))
			assert.True(t, strings.HasSuffix(line, "future_task_test.go:207"))
			//
			line = lines[2]
			assert.Equal(t, "   --- End of error stack trace ---", line)
			//
			line = lines[3]
			assert.True(t, strings.HasPrefix(line, "   created by github.com/timandy/routine.TestFutureTask_Fail_RuntimeError()"))
			assert.True(t, strings.HasSuffix(line, "future_task_test.go:201"))
		}
	}()
	//
	task := NewFutureTask(func(task FutureTask) any { return nil })
	go func() {
		defer func() {
			if cause := recover(); cause != nil {
				task.Fail(NewRuntimeError(cause))
			}
		}()
		panic(1)
	}()
	task.Get()
	assert.Fail(t, "should not be here")
}

func TestFutureTask_Get_Nil(t *testing.T) {
	run := false
	task := NewFutureTask(func(task FutureTask) any { return nil })
	go func() {
		time.Sleep(100 * time.Millisecond)
		run = true
		task.Complete(nil)
	}()
	assert.Nil(t, task.Get())
	assert.True(t, run)
}

func TestFutureTask_Get_Common(t *testing.T) {
	run := false
	task := NewFutureTask(func(task FutureTask) any { return nil })
	go func() {
		time.Sleep(100 * time.Millisecond)
		run = true
		task.Complete(1)
	}()
	assert.Equal(t, 1, task.Get())
	assert.True(t, run)
}

func TestFutureTask_GetWithTimeout_Complete(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	run := false
	task := NewFutureTask(func(task FutureTask) any { return nil })
	go func() {
		defer wg.Done()
		//
		if task.IsCanceled() {
			return
		}
		run = true
		task.Complete(1)
	}()
	assert.Equal(t, 1, task.GetWithTimeout(100*time.Millisecond))
	assert.True(t, run)
	//
	wg.Wait()
}

func TestFutureTask_GetWithTimeout_Fail(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	run := false
	task := NewFutureTask(func(task FutureTask) any { return nil })
	go func() {
		defer wg.Done()
		//
		if task.IsCanceled() {
			return
		}
		run = true
		task.Fail(1)
	}()
	assert.Panics(t, func() {
		task.GetWithTimeout(100 * time.Millisecond)
	})
	assert.True(t, run)
	//
	assert.True(t, task.IsFailed())
	assert.Equal(t, "1", task.(*futureTask).error.Message())
	assert.Nil(t, task.(*futureTask).error.Cause())
	//
	wg.Wait()
}

func TestFutureTask_GetWithTimeout_Timeout(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	run := false
	task := NewFutureTask(func(task FutureTask) any { return nil })
	go func() {
		defer wg.Done()
		//
		time.Sleep(100 * time.Millisecond)
		if task.IsCanceled() {
			return
		}
		run = true
		task.Complete(nil)
	}()
	assert.Panics(t, func() {
		task.GetWithTimeout(1 * time.Millisecond)
	})
	assert.False(t, run)
	//
	assert.True(t, task.IsCanceled())
	assert.Equal(t, "Task execution timeout after 1ms.", task.(*futureTask).error.Message())
	assert.Nil(t, task.(*futureTask).error.Cause())
	//
	wg.Wait()
}

func TestFutureTask_Run_AfterCancel(t *testing.T) {
	run := false
	task := NewFutureTask(func(task FutureTask) any {
		run = true
		return nil
	})
	task.Cancel()
	task.Run()
	assert.Panics(t, func() {
		task.Get()
	})
	assert.True(t, task.IsCanceled())
	assert.False(t, run)
}

func TestFutureTask_Run_AfterFail(t *testing.T) {
	run := false
	task := NewFutureTask(func(task FutureTask) any {
		run = true
		return nil
	})
	task.Fail("failed.")
	task.Run()
	assert.Panics(t, func() {
		task.Get()
	})
	assert.True(t, task.IsFailed())
	assert.False(t, run)
}

func TestFutureTask_Run_AfterComplete(t *testing.T) {
	run := false
	task := NewFutureTask(func(task FutureTask) any {
		run = true
		return nil
	})
	task.Complete(1)
	task.Run()
	assert.Equal(t, 1, task.Get())
	assert.True(t, task.IsDone())
	assert.False(t, run)
}

func TestFutureTask_Run_AfterRun(t *testing.T) {
	var run int32 = 0
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg2 := &sync.WaitGroup{}
	wg2.Add(1)
	task := NewFutureTask(func(task FutureTask) any {
		atomic.AddInt32(&run, 1)
		wg.Done()
		wg2.Wait()
		return 1
	})
	go task.Run()
	wg.Wait()
	task.Run()
	wg2.Done()
	assert.Equal(t, 1, task.Get())
	assert.True(t, task.IsDone())
	assert.Equal(t, int32(1), atomic.LoadInt32(&run))
}

func TestFutureTask_Run_Normal(t *testing.T) {
	run := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	task := NewFutureTask(func(task FutureTask) any {
		run = true
		return 1
	})
	go task.Run()
	assert.Equal(t, 1, task.Get())
	assert.True(t, task.IsDone())
	assert.True(t, run)
}

func TestFutureTask_Run_Error(t *testing.T) {
	run := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	task := NewFutureTask(func(task FutureTask) any {
		run = true
		panic(1)
	})
	go task.Run()
	assert.Panics(t, func() {
		task.Get()
	})
	assert.True(t, task.IsFailed())
	assert.True(t, run)
	//
	defer func() {
		cause := recover()
		assert.NotNil(t, cause)
		assert.Implements(t, (*RuntimeError)(nil), cause)
		err := cause.(RuntimeError)
		assert.Equal(t, "1", err.Message())
		lines := strings.Split(err.Error(), newLine)
		assert.Equal(t, 5, len(lines))
		//
		line := lines[0]
		assert.Equal(t, "RuntimeError: 1", line)
		//
		line = lines[1]
		assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestFutureTask_Run_Error."))
		assert.True(t, strings.HasSuffix(line, "future_task_test.go:397"))
		//
		line = lines[2]
		assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.(*futureTask).Run()"))
		assert.True(t, strings.HasSuffix(line, "future_task.go:108"))
		//
		line = lines[3]
		assert.Equal(t, "   --- End of error stack trace ---", line)
		//
		line = lines[4]
		assert.True(t, strings.HasPrefix(line, "   created by github.com/timandy/routine.TestFutureTask_Run_Error()"))
		assert.True(t, strings.HasSuffix(line, "future_task_test.go:399"))
	}()
	task.Get()
	assert.Fail(t, "should not be here")
}

func TestFutureTask_Run_RuntimeError(t *testing.T) {
	run := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	task := NewFutureTask(func(task FutureTask) any {
		run = true
		err := NewRuntimeError(1)
		panic(err)
	})
	go task.Run()
	assert.Panics(t, func() {
		task.Get()
	})
	assert.True(t, task.IsFailed())
	assert.True(t, run)
	//
	defer func() {
		cause := recover()
		assert.NotNil(t, cause)
		assert.Implements(t, (*RuntimeError)(nil), cause)
		err := cause.(RuntimeError)
		assert.Equal(t, "1", err.Message())
		lines := strings.Split(err.Error(), newLine)
		assert.Equal(t, 5, len(lines))
		//
		line := lines[0]
		assert.Equal(t, "RuntimeError: 1", line)
		//
		line = lines[1]
		assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestFutureTask_Run_RuntimeError."))
		assert.True(t, strings.HasSuffix(line, "future_task_test.go:443"))
		//
		line = lines[2]
		assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.(*futureTask).Run()"))
		assert.True(t, strings.HasSuffix(line, "future_task.go:108"))
		//
		line = lines[3]
		assert.Equal(t, "   --- End of error stack trace ---", line)
		//
		line = lines[4]
		assert.True(t, strings.HasPrefix(line, "   created by github.com/timandy/routine.TestFutureTask_Run_RuntimeError()"))
		assert.True(t, strings.HasSuffix(line, "future_task_test.go:446"))
	}()
	task.Get()
	assert.Fail(t, "should not be here")
}

func TestFutureTask_Routine_Complete(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	task := GoWaitResult(func(token CancelToken) any {
		defer wg.Done()
		//
		if token.IsCanceled() {
			panic("canceled")
		}
		time.Sleep(1 * time.Millisecond)
		return 1
	})
	assert.Equal(t, 1, task.GetWithTimeout(100*time.Millisecond))
	//
	wg.Wait()
}

func TestFutureTask_Routine_Cancel(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	task := GoWaitResult(func(token CancelToken) any {
		defer wg.Done()
		//
		token.Cancel()
		return 1
	})
	assert.Panics(t, func() {
		task.GetWithTimeout(100 * time.Millisecond)
	})
	assert.True(t, task.IsCanceled())
	assert.Equal(t, "Task was canceled.", task.(*futureTask).error.Message())
	assert.Nil(t, task.(*futureTask).error.Cause())
	//
	wg.Wait()
}

func TestFutureTask_Routine_CancelInParent(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg2 := &sync.WaitGroup{}
	wg2.Add(1)
	//
	finished := false
	task := GoWaitResult(func(token CancelToken) any {
		wg2.Done()
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
	wg2.Wait()
	task.Cancel()
	//
	wg.Wait()
	//
	assert.False(t, finished)
	assert.True(t, task.IsCanceled())
	assert.Equal(t, "Task was canceled.", task.(*futureTask).error.Message())
	assert.Nil(t, task.(*futureTask).error.Cause())
}

func TestFutureTask_Routine_Fail(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	task := GoWaitResult(func(token CancelToken) any {
		defer wg.Done()
		//
		if token.IsCanceled() {
			return 1
		}
		panic("something error")
	})
	assert.Panics(t, func() {
		task.GetWithTimeout(100 * time.Millisecond)
	})
	assert.True(t, task.IsFailed())
	assert.Equal(t, "something error", task.(*futureTask).error.Message())
	assert.Nil(t, task.(*futureTask).error.Cause())
	//
	wg.Wait()
}

func TestFutureTask_Routine_Timeout(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	task := GoWaitResult(func(token CancelToken) any {
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
		task.GetWithTimeout(1 * time.Millisecond)
	})
	assert.True(t, task.IsCanceled())
	assert.Equal(t, "Task execution timeout after 1ms.", task.(*futureTask).error.Message())
	assert.Nil(t, task.(*futureTask).error.Cause())
	//
	wg.Wait()
}

func TestFutureTask_Routine_TimeoutThenComplete(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	task := GoWait(func(token CancelToken) {
		defer wg.Done()
		//
		ft := token.(*futureTask)
		ft.result = 1
		assert.True(t, atomic.CompareAndSwapInt32(&ft.state, taskStateRunning, taskStateCompleted))
		time.Sleep(50 * time.Millisecond)
		ft.await.Done()
	})
	assert.Equal(t, 1, task.GetWithTimeout(10*time.Millisecond))
	assert.Equal(t, 1, task.Get())
	//
	wg.Wait()
}

func TestFutureTask_Routine_TimeoutThenCancel(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	task := GoWait(func(token CancelToken) {
		defer wg.Done()
		//
		ft := token.(*futureTask)
		ft.error = NewRuntimeError("canceled.")
		assert.True(t, atomic.CompareAndSwapInt32(&ft.state, taskStateRunning, taskStateCanceled))
		time.Sleep(50 * time.Millisecond)
		ft.await.Done()
	})
	assert.Panics(t, func() {
		task.GetWithTimeout(10 * time.Millisecond)
	})
	//
	assert.True(t, task.IsCanceled())
	assert.Equal(t, "canceled.", task.(*futureTask).error.Message())
	assert.Nil(t, task.(*futureTask).error.Cause())
	assert.Panics(t, func() {
		task.Get()
	})
	//
	wg.Wait()
}

func TestFutureTask_Routine_TimeoutThenFail(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//
	task := GoWait(func(token CancelToken) {
		defer wg.Done()
		//
		ft := token.(*futureTask)
		ft.error = NewRuntimeError("failed.")
		assert.True(t, atomic.CompareAndSwapInt32(&ft.state, taskStateRunning, taskStateFailed))
		time.Sleep(50 * time.Millisecond)
		ft.await.Done()
	})
	assert.Panics(t, func() {
		task.GetWithTimeout(10 * time.Millisecond)
	})
	//
	assert.True(t, task.IsFailed())
	assert.Equal(t, "failed.", task.(*futureTask).error.Message())
	assert.Nil(t, task.(*futureTask).error.Cause())
	assert.Panics(t, func() {
		task.Get()
	})
	//
	wg.Wait()
}

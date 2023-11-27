package routine

import (
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunnable(t *testing.T) {
	count := 0
	var runnable Runnable = func() {
		count++
	}
	runnable()
	assert.Equal(t, 1, count)
	//
	var fun func() = runnable
	fun()
	assert.Equal(t, 2, count)
}

func TestCallable(t *testing.T) {
	var callable Callable[string] = func() string {
		return "Hello"
	}
	assert.Equal(t, "Hello", callable())
	//
	var fun func() string = callable
	assert.Equal(t, "Hello", fun())
}

func TestCancelRunnable(t *testing.T) {
	count := 0
	var cancelRunnable CancelRunnable = func(token CancelToken) {
		count++
	}
	cancelRunnable(nil)
	assert.Equal(t, 1, count)
	//
	var fun func(CancelToken) = cancelRunnable
	fun(nil)
	assert.Equal(t, 2, count)
}

func TestCancelCallable(t *testing.T) {
	var cancelCallable CancelCallable[string] = func(token CancelToken) string {
		return "Hello"
	}
	assert.Equal(t, "Hello", cancelCallable(nil))
	//
	var fun func(CancelToken) string = cancelCallable
	assert.Equal(t, "Hello", fun(nil))
}

func TestWrapTask_NoContext(t *testing.T) {
	run := false
	wrappedRun := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	tls := NewThreadLocal[string]()
	tlsInherit := NewInheritableThreadLocal[string]()
	tls.Set("hello")
	tlsInherit.Set("inherit hello")
	assert.Equal(t, "hello", tls.Get())
	assert.Equal(t, "inherit hello", tlsInherit.Get())
	task := WrapTask(func() {
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, "inherit hello", tlsInherit.Get())
		tls.Set("世界")
		tlsInherit.Set("inherit 世界")
		assert.Equal(t, "世界", tls.Get())
		assert.Equal(t, "inherit 世界", tlsInherit.Get())
		wrappedRun = true
	})
	tls.Set("world")
	tlsInherit.Set("inherit world")
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	go func() {
		task.Run()
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, "", tlsInherit.Get())
		run = true
		wg.Done()
	}()
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	wg.Wait()
	assert.True(t, wrappedRun)
	assert.True(t, run)
}

func TestWrapTask_HasContext(t *testing.T) {
	run := false
	wrappedRun := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	tls := NewThreadLocal[string]()
	tlsInherit := NewInheritableThreadLocal[string]()
	tls.Set("hello")
	tlsInherit.Set("inherit hello")
	assert.Equal(t, "hello", tls.Get())
	assert.Equal(t, "inherit hello", tlsInherit.Get())
	task := WrapTask(func() {
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, "inherit hello", tlsInherit.Get())
		tls.Set("世界")
		tlsInherit.Set("inherit 世界")
		assert.Equal(t, "世界", tls.Get())
		assert.Equal(t, "inherit 世界", tlsInherit.Get())
		wrappedRun = true
	})
	tls.Set("world")
	tlsInherit.Set("inherit world")
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	go func() {
		tls.Set("你好")
		tlsInherit.Set("inherit 你好")
		task.Run()
		assert.Equal(t, "你好", tls.Get())
		assert.Equal(t, "inherit 你好", tlsInherit.Get())
		run = true
		wg.Done()
	}()
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	wg.Wait()
	assert.True(t, wrappedRun)
	assert.True(t, run)
}

func TestWrapTask_Complete_ThenFail(t *testing.T) {
	newStdout, oldStdout := captureStdout()
	defer restoreStdout(newStdout, oldStdout)
	//
	run := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg2 := &sync.WaitGroup{}
	wg2.Add(1)
	wg3 := &sync.WaitGroup{}
	wg3.Add(1)
	task := WrapTask(func() {
		wg.Done()  //1
		wg2.Wait() //4
		run = true
		wg3.Done() //5
		panic(1)
	})
	go task.Run()
	wg.Wait() //2
	task.Complete(nil)
	assert.Nil(t, task.Get())
	wg2.Done() //3
	wg3.Wait() //6
	assert.True(t, task.IsDone())
	assert.False(t, task.IsFailed())
	assert.False(t, task.IsCanceled())
	assert.True(t, run)
	//
	time.Sleep(10 * time.Millisecond)
	output := readAll(newStdout)
	assert.Equal(t, "", output)
}

func TestWrapWaitTask_NoContext(t *testing.T) {
	run := false
	wrappedRun := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	tls := NewThreadLocal[string]()
	tlsInherit := NewInheritableThreadLocal[string]()
	tls.Set("hello")
	tlsInherit.Set("inherit hello")
	assert.Equal(t, "hello", tls.Get())
	assert.Equal(t, "inherit hello", tlsInherit.Get())
	task := WrapWaitTask(func(token CancelToken) {
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, "inherit hello", tlsInherit.Get())
		tls.Set("世界")
		tlsInherit.Set("inherit 世界")
		assert.Equal(t, "世界", tls.Get())
		assert.Equal(t, "inherit 世界", tlsInherit.Get())
		wrappedRun = true
	})
	tls.Set("world")
	tlsInherit.Set("inherit world")
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	go func() {
		task.Run()
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, "", tlsInherit.Get())
		run = true
		wg.Done()
	}()
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	assert.Nil(t, task.Get())
	assert.True(t, wrappedRun)
	wg.Wait()
	assert.True(t, wrappedRun)
	assert.True(t, run)
}

func TestWrapWaitTask_NoContext_Cancel(t *testing.T) {
	run := false
	wrappedRun := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	tls := NewThreadLocal[string]()
	tlsInherit := NewInheritableThreadLocal[string]()
	tls.Set("hello")
	tlsInherit.Set("inherit hello")
	assert.Equal(t, "hello", tls.Get())
	assert.Equal(t, "inherit hello", tlsInherit.Get())
	task := WrapWaitTask(func(token CancelToken) {
		for i := 0; i < 1000; i++ {
			if token.IsCanceled() {
				return
			}
			time.Sleep(1 * time.Millisecond)
		}
		wrappedRun = true
	})
	tls.Set("world")
	tlsInherit.Set("inherit world")
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	go func() {
		task.Run()
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, "", tlsInherit.Get())
		run = true
		wg.Done()
	}()
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	task.Cancel()
	assert.True(t, task.IsCanceled())
	assert.Panics(t, func() {
		task.Get()
	})
	assert.False(t, wrappedRun)
	wg.Wait()
	assert.False(t, wrappedRun)
	assert.True(t, run)
}

func TestWrapWaitTask_HasContext(t *testing.T) {
	run := false
	wrappedRun := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	tls := NewThreadLocal[string]()
	tlsInherit := NewInheritableThreadLocal[string]()
	tls.Set("hello")
	tlsInherit.Set("inherit hello")
	assert.Equal(t, "hello", tls.Get())
	assert.Equal(t, "inherit hello", tlsInherit.Get())
	task := WrapWaitTask(func(token CancelToken) {
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, "inherit hello", tlsInherit.Get())
		tls.Set("世界")
		tlsInherit.Set("inherit 世界")
		assert.Equal(t, "世界", tls.Get())
		assert.Equal(t, "inherit 世界", tlsInherit.Get())
		wrappedRun = true
	})
	tls.Set("world")
	tlsInherit.Set("inherit world")
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	go func() {
		tls.Set("你好")
		tlsInherit.Set("inherit 你好")
		task.Run()
		assert.Equal(t, "你好", tls.Get())
		assert.Equal(t, "inherit 你好", tlsInherit.Get())
		run = true
		wg.Done()
	}()
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	assert.Nil(t, task.Get())
	assert.True(t, wrappedRun)
	wg.Wait()
	assert.True(t, wrappedRun)
	assert.True(t, run)
}

func TestWrapWaitTask_HasContext_Cancel(t *testing.T) {
	run := false
	wrappedRun := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	tls := NewThreadLocal[string]()
	tlsInherit := NewInheritableThreadLocal[string]()
	tls.Set("hello")
	tlsInherit.Set("inherit hello")
	assert.Equal(t, "hello", tls.Get())
	assert.Equal(t, "inherit hello", tlsInherit.Get())
	task := WrapWaitTask(func(token CancelToken) {
		for i := 0; i < 1000; i++ {
			if token.IsCanceled() {
				return
			}
			time.Sleep(1 * time.Millisecond)
		}
		wrappedRun = true
	})
	tls.Set("world")
	tlsInherit.Set("inherit world")
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	go func() {
		tls.Set("你好")
		tlsInherit.Set("inherit 你好")
		task.Run()
		assert.Equal(t, "你好", tls.Get())
		assert.Equal(t, "inherit 你好", tlsInherit.Get())
		run = true
		wg.Done()
	}()
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	task.Cancel()
	assert.True(t, task.IsCanceled())
	assert.Panics(t, func() {
		task.Get()
	})
	assert.False(t, wrappedRun)
	wg.Wait()
	assert.False(t, wrappedRun)
	assert.True(t, run)
}

func TestWrapWaitTask_Complete_ThenFail(t *testing.T) {
	newStdout, oldStdout := captureStdout()
	defer restoreStdout(newStdout, oldStdout)
	//
	run := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg2 := &sync.WaitGroup{}
	wg2.Add(1)
	wg3 := &sync.WaitGroup{}
	wg3.Add(1)
	task := WrapWaitTask(func(token CancelToken) {
		wg.Done()  //1
		wg2.Wait() //4
		run = true
		wg3.Done() //5
		panic(1)
	})
	go task.Run()
	wg.Wait() //2
	task.Complete(nil)
	assert.Nil(t, task.Get())
	wg2.Done() //3
	wg3.Wait() //6
	assert.True(t, task.IsDone())
	assert.False(t, task.IsFailed())
	assert.False(t, task.IsCanceled())
	assert.True(t, run)
	//
	time.Sleep(10 * time.Millisecond)
	output := readAll(newStdout)
	assert.Equal(t, "", output)
}

func TestWrapWaitResultTask_NoContext(t *testing.T) {
	run := false
	wrappedRun := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	tls := NewThreadLocal[string]()
	tlsInherit := NewInheritableThreadLocal[string]()
	tls.Set("hello")
	tlsInherit.Set("inherit hello")
	assert.Equal(t, "hello", tls.Get())
	assert.Equal(t, "inherit hello", tlsInherit.Get())
	task := WrapWaitResultTask(func(token CancelToken) any {
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, "inherit hello", tlsInherit.Get())
		tls.Set("世界")
		tlsInherit.Set("inherit 世界")
		assert.Equal(t, "世界", tls.Get())
		assert.Equal(t, "inherit 世界", tlsInherit.Get())
		wrappedRun = true
		return 1
	})
	tls.Set("world")
	tlsInherit.Set("inherit world")
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	go func() {
		task.Run()
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, "", tlsInherit.Get())
		run = true
		wg.Done()
	}()
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	assert.Equal(t, 1, task.Get())
	assert.True(t, wrappedRun)
	wg.Wait()
	assert.True(t, wrappedRun)
	assert.True(t, run)
}

func TestWrapWaitResultTask_NoContext_Cancel(t *testing.T) {
	run := false
	wrappedRun := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	tls := NewThreadLocal[string]()
	tlsInherit := NewInheritableThreadLocal[string]()
	tls.Set("hello")
	tlsInherit.Set("inherit hello")
	assert.Equal(t, "hello", tls.Get())
	assert.Equal(t, "inherit hello", tlsInherit.Get())
	task := WrapWaitResultTask(func(token CancelToken) any {
		for i := 0; i < 1000; i++ {
			if token.IsCanceled() {
				return 0
			}
			time.Sleep(1 * time.Millisecond)
		}
		wrappedRun = true
		return 1
	})
	tls.Set("world")
	tlsInherit.Set("inherit world")
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	go func() {
		task.Run()
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, "", tlsInherit.Get())
		run = true
		wg.Done()
	}()
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	task.Cancel()
	assert.True(t, task.IsCanceled())
	assert.Panics(t, func() {
		task.Get()
	})
	assert.False(t, wrappedRun)
	wg.Wait()
	assert.False(t, wrappedRun)
	assert.True(t, run)
}

func TestWrapWaitResultTask_HasContext(t *testing.T) {
	run := false
	wrappedRun := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	tls := NewThreadLocal[string]()
	tlsInherit := NewInheritableThreadLocal[string]()
	tls.Set("hello")
	tlsInherit.Set("inherit hello")
	assert.Equal(t, "hello", tls.Get())
	assert.Equal(t, "inherit hello", tlsInherit.Get())
	task := WrapWaitResultTask(func(token CancelToken) any {
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, "inherit hello", tlsInherit.Get())
		tls.Set("世界")
		tlsInherit.Set("inherit 世界")
		assert.Equal(t, "世界", tls.Get())
		assert.Equal(t, "inherit 世界", tlsInherit.Get())
		wrappedRun = true
		return 1
	})
	tls.Set("world")
	tlsInherit.Set("inherit world")
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	go func() {
		tls.Set("你好")
		tlsInherit.Set("inherit 你好")
		task.Run()
		assert.Equal(t, "你好", tls.Get())
		assert.Equal(t, "inherit 你好", tlsInherit.Get())
		run = true
		wg.Done()
	}()
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	assert.Equal(t, 1, task.Get())
	assert.True(t, wrappedRun)
	wg.Wait()
	assert.True(t, wrappedRun)
	assert.True(t, run)
}

func TestWrapWaitResultTask_HasContext_Cancel(t *testing.T) {
	run := false
	wrappedRun := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	tls := NewThreadLocal[string]()
	tlsInherit := NewInheritableThreadLocal[string]()
	tls.Set("hello")
	tlsInherit.Set("inherit hello")
	assert.Equal(t, "hello", tls.Get())
	assert.Equal(t, "inherit hello", tlsInherit.Get())
	task := WrapWaitResultTask(func(token CancelToken) any {
		for i := 0; i < 1000; i++ {
			if token.IsCanceled() {
				return 0
			}
			time.Sleep(1 * time.Millisecond)
		}
		wrappedRun = true
		return 1
	})
	tls.Set("world")
	tlsInherit.Set("inherit world")
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	go func() {
		tls.Set("你好")
		tlsInherit.Set("inherit 你好")
		task.Run()
		assert.Equal(t, "你好", tls.Get())
		assert.Equal(t, "inherit 你好", tlsInherit.Get())
		run = true
		wg.Done()
	}()
	assert.Equal(t, "world", tls.Get())
	assert.Equal(t, "inherit world", tlsInherit.Get())
	task.Cancel()
	assert.True(t, task.IsCanceled())
	assert.Panics(t, func() {
		task.Get()
	})
	assert.False(t, wrappedRun)
	wg.Wait()
	assert.False(t, wrappedRun)
	assert.True(t, run)
}

func TestWrapWaitResultTask_Complete_ThenFail(t *testing.T) {
	newStdout, oldStdout := captureStdout()
	defer restoreStdout(newStdout, oldStdout)
	//
	run := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg2 := &sync.WaitGroup{}
	wg2.Add(1)
	wg3 := &sync.WaitGroup{}
	wg3.Add(1)
	task := WrapWaitResultTask(func(token CancelToken) any {
		wg.Done()  //1
		wg2.Wait() //4
		run = true
		wg3.Done() //5
		panic(1)
	})
	go task.Run()
	wg.Wait() //2
	task.Complete(nil)
	assert.Nil(t, task.Get())
	wg2.Done() //3
	wg3.Wait() //6
	assert.True(t, task.IsDone())
	assert.False(t, task.IsFailed())
	assert.False(t, task.IsCanceled())
	assert.True(t, run)
	//
	time.Sleep(10 * time.Millisecond)
	output := readAll(newStdout)
	assert.Equal(t, "", output)
}

func TestGo_Error(t *testing.T) {
	newStdout, oldStdout := captureStdout()
	defer restoreStdout(newStdout, oldStdout)
	//
	run := false
	assert.NotPanics(t, func() {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		Go(func() {
			run = true
			wg.Done()
			panic("error")
		})
		wg.Wait()
	})
	assert.True(t, run)
	//
	time.Sleep(10 * time.Millisecond)
	output := readAll(newStdout)
	lines := strings.Split(output, newLine)
	assert.Equal(t, 7, len(lines))
	//
	line := lines[0]
	assert.Equal(t, "RuntimeError: error", line)
	//
	line = lines[1]
	assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestGo_Error."))
	assert.True(t, strings.HasSuffix(line, "api_routine_test.go:601"))
	//
	line = lines[2]
	assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.inheritedTask.run()"))
	assert.True(t, strings.HasSuffix(line, "routine.go:31"))
	//
	line = lines[3]
	assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.(*futureTask[...]).Run()"))
	assert.True(t, strings.HasSuffix(line, "future_task.go:108"))
	//
	line = lines[4]
	assert.Equal(t, "   --- End of error stack trace ---", line)
	//
	line = lines[5]
	assert.True(t, strings.HasPrefix(line, "   created by github.com/timandy/routine.Go()"))
	assert.True(t, strings.HasSuffix(line, "api_routine.go:49"))
	//
	line = lines[6]
	assert.Equal(t, "", line)
}

func TestGo_Nil(t *testing.T) {
	assert.Nil(t, createInheritedMap())
	//
	run := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	Go(func() {
		assert.Nil(t, createInheritedMap())
		run = true
		wg.Done()
	})
	wg.Wait()
	assert.True(t, run)
}

func TestGo_Value(t *testing.T) {
	tls := NewThreadLocal[string]()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	inheritableTls := NewInheritableThreadLocal[string]()
	inheritableTls.Set("World")
	assert.Equal(t, "World", inheritableTls.Get())
	//
	assert.NotNil(t, createInheritedMap())
	//
	run := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	Go(func() {
		assert.NotNil(t, createInheritedMap())
		//
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, "World", inheritableTls.Get())
		//
		tls.Set("Hello2")
		assert.Equal(t, "Hello2", tls.Get())
		//
		inheritableTls.Remove()
		assert.Equal(t, "", inheritableTls.Get())
		//
		run = true
		wg.Done()
	})
	wg.Wait()
	assert.True(t, run)
	//
	assert.Equal(t, "Hello", tls.Get())
	assert.Equal(t, "World", inheritableTls.Get())
}

func TestGo_Cross(t *testing.T) {
	tls := NewThreadLocal[string]()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	wg := &sync.WaitGroup{}
	wg.Add(1)
	Go(func() {
		assert.Equal(t, "", tls.Get())
		wg.Done()
	})
	wg.Wait()
}

func TestGoWait_Error(t *testing.T) {
	run := false
	task := GoWait(func(token CancelToken) {
		run = true
		panic("error")
	})
	assert.Panics(t, func() {
		task.Get()
	})
	assert.True(t, run)
	//
	defer func() {
		cause := recover()
		assert.NotNil(t, cause)
		assert.Implements(t, (*RuntimeError)(nil), cause)
		err := cause.(RuntimeError)
		lines := strings.Split(err.Error(), newLine)
		assert.Equal(t, 6, len(lines))
		//
		line := lines[0]
		assert.Equal(t, "RuntimeError: error", line)
		//
		line = lines[1]
		assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestGoWait_Error."))
		assert.True(t, strings.HasSuffix(line, "api_routine_test.go:707"))
		//
		line = lines[2]
		assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.inheritedWaitTask.run()"))
		assert.True(t, strings.HasSuffix(line, "routine.go:70"))
		//
		line = lines[3]
		assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.(*futureTask[...]).Run()"))
		assert.True(t, strings.HasSuffix(line, "future_task.go:108"))
		//
		line = lines[4]
		assert.Equal(t, "   --- End of error stack trace ---", line)
		//
		line = lines[5]
		assert.True(t, strings.HasPrefix(line, "   created by github.com/timandy/routine.GoWait()"))
		assert.True(t, strings.HasSuffix(line, "api_routine.go:57"))
	}()
	task.Get()
}

func TestGoWait_Nil(t *testing.T) {
	assert.Nil(t, createInheritedMap())
	//
	run := false
	task := GoWait(func(token CancelToken) {
		assert.Nil(t, createInheritedMap())
		run = true
	})
	assert.Nil(t, task.Get())
	assert.True(t, run)
}

func TestGoWait_Value(t *testing.T) {
	tls := NewThreadLocal[string]()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	inheritableTls := NewInheritableThreadLocal[string]()
	inheritableTls.Set("World")
	assert.Equal(t, "World", inheritableTls.Get())
	//
	assert.NotNil(t, createInheritedMap())
	//
	run := false
	task := GoWait(func(token CancelToken) {
		assert.NotNil(t, createInheritedMap())
		//
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, "World", inheritableTls.Get())
		//
		tls.Set("Hello2")
		assert.Equal(t, "Hello2", tls.Get())
		//
		inheritableTls.Remove()
		assert.Equal(t, "", inheritableTls.Get())
		//
		run = true
	})
	assert.Nil(t, task.Get())
	assert.True(t, run)
	//
	assert.Equal(t, "Hello", tls.Get())
	assert.Equal(t, "World", inheritableTls.Get())
}

func TestGoWait_Cross(t *testing.T) {
	tls := NewThreadLocal[string]()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	GoWait(func(token CancelToken) {
		assert.Equal(t, "", tls.Get())
	}).Get()
}

func TestGoWaitResult_Error(t *testing.T) {
	run := false
	task := GoWaitResult(func(token CancelToken) int {
		run = true
		if run {
			panic("error")
		}
		return 1
	})
	assert.Panics(t, func() {
		task.Get()
	})
	assert.True(t, run)
	//
	defer func() {
		cause := recover()
		assert.NotNil(t, cause)
		assert.Implements(t, (*RuntimeError)(nil), cause)
		err := cause.(RuntimeError)
		lines := strings.Split(err.Error(), newLine)
		assert.True(t, len(lines) == 6 || len(lines) == 7)
		//
		line := lines[0]
		assert.Equal(t, "RuntimeError: error", line)
		//
		line = lines[1]
		assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.TestGoWaitResult_Error."))
		assert.True(t, strings.HasSuffix(line, "api_routine_test.go:807"))
		//
		line = lines[2]
		assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.inheritedWaitResultTask[...].run()"))
		assert.True(t, strings.HasSuffix(line, "routine.go:109"))
		//
		lineOffset := 0
		if len(lines) == 7 {
			line = lines[3+lineOffset]
			lineOffset = 1
			assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.WrapWaitResultTask[...].func1()"))
			assert.True(t, strings.HasSuffix(line, "api_routine.go:41"))
		}
		//
		line = lines[3+lineOffset]
		assert.True(t, strings.HasPrefix(line, "   at github.com/timandy/routine.(*futureTask[...]).Run()"))
		assert.True(t, strings.HasSuffix(line, "future_task.go:108"))
		//
		line = lines[4+lineOffset]
		assert.Equal(t, "   --- End of error stack trace ---", line)
		//
		line = lines[5+lineOffset]
		assert.True(t, strings.HasPrefix(line, "   created by github.com/timandy/routine.GoWaitResult[...]()"))
		assert.True(t, strings.HasSuffix(line, "api_routine.go:66"))
	}()
	task.Get()
}

func TestGoWaitResult_Nil(t *testing.T) {
	assert.Nil(t, createInheritedMap())
	//
	run := false
	task := GoWaitResult(func(token CancelToken) bool {
		assert.Nil(t, createInheritedMap())
		run = true
		return true
	})
	assert.True(t, task.Get())
	assert.True(t, run)
}

func TestGoWaitResult_Value(t *testing.T) {
	tls := NewThreadLocal[string]()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	inheritableTls := NewInheritableThreadLocal[string]()
	inheritableTls.Set("World")
	assert.Equal(t, "World", inheritableTls.Get())
	//
	assert.NotNil(t, createInheritedMap())
	//
	run := false
	task := GoWaitResult(func(token CancelToken) bool {
		assert.NotNil(t, createInheritedMap())
		//
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, "World", inheritableTls.Get())
		//
		tls.Set("Hello2")
		assert.Equal(t, "Hello2", tls.Get())
		//
		inheritableTls.Remove()
		assert.Equal(t, "", inheritableTls.Get())
		//
		run = true
		return true
	})
	assert.True(t, task.Get())
	assert.True(t, run)
	//
	assert.Equal(t, "Hello", tls.Get())
	assert.Equal(t, "World", inheritableTls.Get())
}

func TestGoWaitResult_Cross(t *testing.T) {
	tls := NewThreadLocal[string]()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	result := GoWaitResult(func(token CancelToken) string {
		assert.Equal(t, "", tls.Get())
		return tls.Get()
	}).Get()
	assert.Equal(t, "", result)
}

func captureStdout() (newStdout, oldStdout *os.File) {
	oldStdout = os.Stdout
	fileName := path.Join(os.TempDir(), "go_test_"+strconv.FormatInt(time.Now().UnixNano(), 10)+".txt")
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	os.Stdout = file
	newStdout = file
	return
}

func restoreStdout(newStdout, oldStdout *os.File) {
	os.Stdout = oldStdout
	if err := newStdout.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(newStdout.Name()); err != nil {
		panic(err)
	}
}

func readAll(rs io.ReadSeeker) string {
	if _, err := rs.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}
	b := make([]byte, 0, 512)
	for {
		if len(b) == cap(b) {
			b = append(b, 0)[:len(b)]
		}
		n, err := rs.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				return string(b)
			}
			panic(err)
		}
	}
}

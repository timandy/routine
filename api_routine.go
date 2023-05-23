package routine

import "fmt"

// Runnable provides a function without return values.
type Runnable func()

// CancelRunnable provides a cancellable function without return values.
type CancelRunnable func(token CancelToken)

// CancelCallable provides a cancellable function that returns a value of type any.
type CancelCallable func(token CancelToken) any

// Go starts a new goroutine, and copy inheritableThreadLocals from current goroutine.
// This function will auto invoke the fun and print error stack when panic occur in goroutine.
func Go(fun Runnable) {
	// backup
	copied := createInheritedMap()
	go func() {
		// catch
		defer func() {
			if cause := recover(); cause != nil {
				fmt.Println(NewRuntimeError(cause).Error())
			}
		}()
		// restore
		t := currentThread(copied != nil)
		if t == nil {
			//copied is nil
			defer func() {
				t = currentThread(false)
				if t != nil {
					t.threadLocals = nil
					t.inheritableThreadLocals = nil
				}
			}()
			fun()
		} else {
			backup := t.inheritableThreadLocals
			defer func() {
				t.threadLocals = nil
				t.inheritableThreadLocals = backup
			}()
			t.threadLocals = nil
			t.inheritableThreadLocals = copied
			fun()
		}
	}()
}

// GoWait starts a new goroutine, and copy inheritableThreadLocals from current goroutine.
// This function return a FutureTask pointer, so we can wait by FutureTask.Get or FutureTask.GetWithTimeout method.
// If panic occur in goroutine, The panic will be trigger again when calling FutureTask.Get or FutureTask.GetWithTimeout method.
func GoWait(fun CancelRunnable) FutureTask {
	task := NewFutureTask()
	// backup
	copied := createInheritedMap()
	go func() {
		// catch
		defer func() {
			if cause := recover(); cause != nil {
				task.Fail(NewRuntimeError(cause))
			}
		}()
		// restore
		t := currentThread(copied != nil)
		if t == nil {
			//copied is nil
			defer func() {
				t = currentThread(false)
				if t != nil {
					t.threadLocals = nil
					t.inheritableThreadLocals = nil
				}
			}()
			fun(task)
			task.Complete(nil)
		} else {
			backup := t.inheritableThreadLocals
			defer func() {
				t.threadLocals = nil
				t.inheritableThreadLocals = backup
			}()
			t.threadLocals = nil
			t.inheritableThreadLocals = copied
			fun(task)
			task.Complete(nil)
		}
	}()
	return task
}

// GoWaitResult starts a new goroutine, and copy inheritableThreadLocals from current goroutine.
// This function return a FutureTask pointer, so we can wait and get result by FutureTask.Get or FutureTask.GetWithTimeout method.
// If panic occur in goroutine, The panic will be trigger again when calling FutureTask.Get or FutureTask.GetWithTimeout method.
func GoWaitResult(fun CancelCallable) FutureTask {
	task := NewFutureTask()
	// backup
	copied := createInheritedMap()
	go func() {
		// catch
		defer func() {
			if cause := recover(); cause != nil {
				task.Fail(NewRuntimeError(cause))
			}
		}()
		// restore
		t := currentThread(copied != nil)
		if t == nil {
			//copied is nil
			defer func() {
				t = currentThread(false)
				if t != nil {
					t.threadLocals = nil
					t.inheritableThreadLocals = nil
				}
			}()
			task.Complete(fun(task))
		} else {
			backup := t.inheritableThreadLocals
			defer func() {
				t.threadLocals = nil
				t.inheritableThreadLocals = backup
			}()
			t.threadLocals = nil
			t.inheritableThreadLocals = copied
			task.Complete(fun(task))
		}
	}()
	return task
}

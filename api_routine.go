package routine

import "fmt"

// Runnable provides a function without return values.
type Runnable func()

// Callable provides a function that returns a value of type Any.
type Callable func() Any

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
// This function return a Future pointer, so we can wait by Future.Get or Future.GetWithTimeout method.
// If panic occur in goroutine, The panic will be trigger again when calling Future.Get or Future.GetWithTimeout method.
func GoWait(fun Runnable) Future {
	fut := NewFuture()
	// backup
	copied := createInheritedMap()
	go func() {
		// catch
		defer func() {
			if cause := recover(); cause != nil {
				fut.Fail(NewRuntimeError(cause))
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
			fut.Complete(nil)
		} else {
			backup := t.inheritableThreadLocals
			defer func() {
				t.threadLocals = nil
				t.inheritableThreadLocals = backup
			}()
			t.threadLocals = nil
			t.inheritableThreadLocals = copied
			fun()
			fut.Complete(nil)
		}
	}()
	return fut
}

// GoWaitResult starts a new goroutine, and copy inheritableThreadLocals from current goroutine.
// This function return a Future pointer, so we can wait and get result by Future.Get or Future.GetWithTimeout method.
// If panic occur in goroutine, The panic will be trigger again when calling Future.Get or Future.GetWithTimeout method.
func GoWaitResult(fun Callable) Future {
	fut := NewFuture()
	// backup
	copied := createInheritedMap()
	go func() {
		// catch
		defer func() {
			if cause := recover(); cause != nil {
				fut.Fail(NewRuntimeError(cause))
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
			fut.Complete(fun())
		} else {
			backup := t.inheritableThreadLocals
			defer func() {
				t.threadLocals = nil
				t.inheritableThreadLocals = backup
			}()
			t.threadLocals = nil
			t.inheritableThreadLocals = copied
			fut.Complete(fun())
		}
	}()
	return fut
}

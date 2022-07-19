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
// This function return a Future pointer, so we can wait by Future.Get method.
// If panic occur in goroutine, The panic will be trigger again when calling Future.Get method.
func GoWait(fun Runnable) Future {
	fea := NewFuture()
	// backup
	copied := createInheritedMap()
	go func() {
		// catch
		defer func() {
			if cause := recover(); cause != nil {
				fea.CompleteError(NewRuntimeError(cause))
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
			fea.Complete(nil)
		} else {
			backup := t.inheritableThreadLocals
			defer func() {
				t.threadLocals = nil
				t.inheritableThreadLocals = backup
			}()
			t.threadLocals = nil
			t.inheritableThreadLocals = copied
			fun()
			fea.Complete(nil)
		}
	}()
	return fea
}

// GoWaitResult starts a new goroutine, and copy inheritableThreadLocals from current goroutine.
// This function return a Future pointer, so we can wait and get result by Future.Get method.
// If panic occur in goroutine, The panic will be trigger again when calling Future.Get method.
func GoWaitResult(fun Callable) Future {
	fea := NewFuture()
	// backup
	copied := createInheritedMap()
	go func() {
		// catch
		defer func() {
			if cause := recover(); cause != nil {
				fea.CompleteError(NewRuntimeError(cause))
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
			fea.Complete(fun())
		} else {
			backup := t.inheritableThreadLocals
			defer func() {
				t.threadLocals = nil
				t.inheritableThreadLocals = backup
			}()
			t.threadLocals = nil
			t.inheritableThreadLocals = copied
			fea.Complete(fun())
		}
	}()
	return fea
}

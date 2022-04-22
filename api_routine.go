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
			if err := recover(); err != nil {
				fmt.Println(NewStackError(err).Error())
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
// This function return a Feature pointer, so we can wait by Feature.Get method.
// If panic occur in goroutine, The panic will be trigger again when calling Feature.Get method.
func GoWait(fun Runnable) Feature {
	fea := NewFeature()
	// backup
	copied := createInheritedMap()
	go func() {
		// catch
		defer func() {
			if err := recover(); err != nil {
				fea.CompleteError(err)
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
// This function return a Feature pointer, so we can wait and get result by Feature.Get method.
// If panic occur in goroutine, The panic will be trigger again when calling Feature.Get method.
func GoWaitResult(fun Callable) Feature {
	fea := NewFeature()
	// backup
	copied := createInheritedMap()
	go func() {
		// catch
		defer func() {
			if err := recover(); err != nil {
				fea.CompleteError(err)
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

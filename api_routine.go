package routine

import "fmt"

// Go starts a new goroutine, and copy all local table from current goroutine.
func Go(fun func()) {
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

// GoWait starts a new goroutine, and copy all local table from current goroutine.
func GoWait(fun func()) Feature {
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

// GoWaitResult starts a new goroutine, and copy all local table from current goroutine.
func GoWaitResult(fun func() Any) Feature {
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

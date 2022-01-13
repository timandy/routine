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
				fe := &StackError{error: err, stackTrace: readStackBuf()}
				fmt.Println(fe.Error())
			}
		}()
		// restore
		t := currentThread(copied != nil)
		if t == nil {
			fun()
		} else {
			backup := t.inheritableThreadLocals
			t.inheritableThreadLocals = copied
			fun()
			t.inheritableThreadLocals = backup
		}
	}()
}

// Go starts a new goroutine, and copy all local table from current goroutine.
func GoWait(fun func()) *Feature {
	fea := feature()
	// backup
	copied := createInheritedMap()
	go func() {
		// catch
		defer func() {
			if err := recover(); err != nil {
				fea.completeError(stackError(err))
			}
		}()
		// restore
		t := currentThread(copied != nil)
		if t == nil {
			fun()
			fea.complete(nil)
		} else {
			backup := t.inheritableThreadLocals
			t.inheritableThreadLocals = copied
			fun()
			fea.complete(nil)
			t.inheritableThreadLocals = backup
		}
	}()
	return fea
}

// Go starts a new goroutine, and copy all local table from current goroutine.
func GoWaitResult(fun func() Any) *Feature {
	fea := feature()
	// backup
	copied := createInheritedMap()
	go func() {
		// catch
		defer func() {
			if err := recover(); err != nil {
				fea.completeError(stackError(err))
			}
		}()
		// restore
		t := currentThread(copied != nil)
		if t == nil {
			fea.complete(fun())
		} else {
			backup := t.inheritableThreadLocals
			t.inheritableThreadLocals = copied
			fea.complete(fun())
			t.inheritableThreadLocals = backup
		}
	}()
	return fea
}

package routine

import "fmt"

type inherited struct {
}

//go:norace
func (inherited) reset() {
	t := currentThread(false)
	if t != nil {
		t.threadLocals = nil
		t.inheritableThreadLocals = nil
	}
}

//go:norace
func (inherited) restore(t *thread, threadLocalsBackup, inheritableThreadLocalsBackup *threadLocalMap) {
	t.threadLocals = threadLocalsBackup
	t.inheritableThreadLocals = inheritableThreadLocalsBackup
}

type inheritedTask struct {
	inherited
	context  *threadLocalMap
	function Runnable
}

//go:norace
func (it inheritedTask) run(task FutureTask[any]) any {
	// catch
	defer func() {
		if cause := recover(); cause != nil {
			task.Fail(cause)
			if err := task.(*futureTask[any]).error; err != nil {
				fmt.Println(err.Error())
			}
		}
	}()
	// restore
	t := currentThread(it.context != nil)
	if t == nil {
		//copied is nil
		defer it.reset()
		it.function()
		return nil
	} else {
		threadLocalsBackup := t.threadLocals
		inheritableThreadLocalsBackup := t.inheritableThreadLocals
		defer it.restore(t, threadLocalsBackup, inheritableThreadLocalsBackup)
		t.threadLocals = nil
		t.inheritableThreadLocals = it.context
		it.function()
		return nil
	}
}

type inheritedWaitTask struct {
	inherited
	context  *threadLocalMap
	function CancelRunnable
}

//go:norace
func (iwt inheritedWaitTask) run(task FutureTask[any]) any {
	// catch
	defer func() {
		if cause := recover(); cause != nil {
			task.Fail(cause)
		}
	}()
	// restore
	t := currentThread(iwt.context != nil)
	if t == nil {
		//copied is nil
		defer iwt.reset()
		iwt.function(task)
		return nil
	} else {
		threadLocalsBackup := t.threadLocals
		inheritableThreadLocalsBackup := t.inheritableThreadLocals
		defer iwt.restore(t, threadLocalsBackup, inheritableThreadLocalsBackup)
		t.threadLocals = nil
		t.inheritableThreadLocals = iwt.context
		iwt.function(task)
		return nil
	}
}

type inheritedWaitResultTask[TResult any] struct {
	inherited
	context  *threadLocalMap
	function CancelCallable[TResult]
}

//go:norace
func (iwrt inheritedWaitResultTask[TResult]) run(task FutureTask[TResult]) TResult {
	// catch
	defer func() {
		if cause := recover(); cause != nil {
			task.Fail(cause)
		}
	}()
	// restore
	t := currentThread(iwrt.context != nil)
	if t == nil {
		//copied is nil
		defer iwrt.reset()
		return iwrt.function(task)
	} else {
		threadLocalsBackup := t.threadLocals
		inheritableThreadLocalsBackup := t.inheritableThreadLocals
		defer iwrt.restore(t, threadLocalsBackup, inheritableThreadLocalsBackup)
		t.threadLocals = nil
		t.inheritableThreadLocals = iwrt.context
		return iwrt.function(task)
	}
}

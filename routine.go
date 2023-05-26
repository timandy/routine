package routine

import "fmt"

type inheritedTask struct {
	context  *threadLocalMap
	function Runnable
}

func (it inheritedTask) run(task FutureTask) any {
	// catch
	defer func() {
		if cause := recover(); cause != nil {
			task.Fail(cause)
			if err := task.(*futureTask).error; err != nil {
				fmt.Println(err.Error())
			}
		}
	}()
	// restore
	t := currentThread(it.context != nil)
	if t == nil {
		//copied is nil
		defer func() {
			t = currentThread(false)
			if t != nil {
				t.threadLocals = nil
				t.inheritableThreadLocals = nil
			}
		}()
		it.function()
		return nil
	} else {
		threadLocalsBackup := t.threadLocals
		inheritableThreadLocalsBackup := t.inheritableThreadLocals
		defer func() {
			t.threadLocals = threadLocalsBackup
			t.inheritableThreadLocals = inheritableThreadLocalsBackup
		}()
		t.threadLocals = nil
		t.inheritableThreadLocals = it.context
		it.function()
		return nil
	}
}

type inheritedWaitTask struct {
	context  *threadLocalMap
	function CancelRunnable
}

func (iwt inheritedWaitTask) run(task FutureTask) any {
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
		defer func() {
			t = currentThread(false)
			if t != nil {
				t.threadLocals = nil
				t.inheritableThreadLocals = nil
			}
		}()
		iwt.function(task)
		return nil
	} else {
		threadLocalsBackup := t.threadLocals
		inheritableThreadLocalsBackup := t.inheritableThreadLocals
		defer func() {
			t.threadLocals = threadLocalsBackup
			t.inheritableThreadLocals = inheritableThreadLocalsBackup
		}()
		t.threadLocals = nil
		t.inheritableThreadLocals = iwt.context
		iwt.function(task)
		return nil
	}
}

type inheritedWaitResultTask struct {
	context  *threadLocalMap
	function CancelCallable
}

func (iwrt inheritedWaitResultTask) run(task FutureTask) any {
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
		defer func() {
			t = currentThread(false)
			if t != nil {
				t.threadLocals = nil
				t.inheritableThreadLocals = nil
			}
		}()
		return iwrt.function(task)
	} else {
		threadLocalsBackup := t.threadLocals
		inheritableThreadLocalsBackup := t.inheritableThreadLocals
		defer func() {
			t.threadLocals = threadLocalsBackup
			t.inheritableThreadLocals = inheritableThreadLocalsBackup
		}()
		t.threadLocals = nil
		t.inheritableThreadLocals = iwrt.context
		return iwrt.function(task)
	}
}

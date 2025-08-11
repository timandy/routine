package routine

import "fmt"

type inheritedTask struct {
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
	defer restoreInheritedMap(it.context)()
	// exec
	it.function()
	return nil
}

type inheritedWaitTask struct {
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
	defer restoreInheritedMap(iwt.context)()
	// exec
	iwt.function(task)
	return nil
}

type inheritedWaitResultTask[TResult any] struct {
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
	defer restoreInheritedMap(iwrt.context)()
	// exec
	return iwrt.function(task)
}

package routine

// Runnable provides a function without return values.
type Runnable func()

// Callable provides a function that returns a value of type any.
type Callable func() any

// CancelRunnable provides a cancellable function without return values.
type CancelRunnable func(token CancelToken)

// CancelCallable provides a cancellable function that returns a value of type any.
type CancelCallable func(token CancelToken) any

// WrapTask create a new task, and capture inheritableThreadLocals from current goroutine.
// This function return a FutureTask instance, so we can wait and get result by FutureTask.Get or FutureTask.GetWithTimeout method.
// This function will not invoke the func. When the returned task run it will print error stack when panic occur.
func WrapTask(fun Runnable) FutureTask {
	ctx := createInheritedMap()
	callable := inheritedTask{context: ctx, function: fun}.run
	return NewFutureTask(callable)
}

// WrapWaitTask create a new task, and capture inheritableThreadLocals from current goroutine.
// This function return a FutureTask instance, so we can wait by FutureTask.Get or FutureTask.GetWithTimeout method.
// This function will not invoke the func. When the returned task run any panic will be caught, The panic will be trigger again when calling FutureTask.Get or FutureTask.GetWithTimeout method.
func WrapWaitTask(fun CancelRunnable) FutureTask {
	ctx := createInheritedMap()
	callable := inheritedWaitTask{context: ctx, function: fun}.run
	return NewFutureTask(callable)
}

// WrapWaitResultTask create a new task, and capture inheritableThreadLocals from current goroutine.
// This function return a FutureTask instance, so we can wait and get result by FutureTask.Get or FutureTask.GetWithTimeout method.
// This function will not invoke the func. When the returned task run any panic will be caught, The panic will be trigger again when calling FutureTask.Get or FutureTask.GetWithTimeout method.
func WrapWaitResultTask(fun CancelCallable) FutureTask {
	ctx := createInheritedMap()
	callable := inheritedWaitResultTask{context: ctx, function: fun}.run
	return NewFutureTask(callable)
}

// Go starts a new goroutine, and copy inheritableThreadLocals from current goroutine.
// This function will auto invoke the func and print error stack when panic occur in goroutine.
func Go(fun Runnable) {
	task := WrapTask(fun)
	go task.Run()
}

// GoWait starts a new goroutine, and copy inheritableThreadLocals from current goroutine.
// This function will auto invoke the func and return a FutureTask instance, so we can wait by FutureTask.Get or FutureTask.GetWithTimeout method.
// If panic occur in goroutine, The panic will be trigger again when calling FutureTask.Get or FutureTask.GetWithTimeout method.
func GoWait(fun CancelRunnable) FutureTask {
	task := WrapWaitTask(fun)
	go task.Run()
	return task
}

// GoWaitResult starts a new goroutine, and copy inheritableThreadLocals from current goroutine.
// This function will auto invoke the func and return a FutureTask instance, so we can wait and get result by FutureTask.Get or FutureTask.GetWithTimeout method.
// If panic occur in goroutine, The panic will be trigger again when calling FutureTask.Get or FutureTask.GetWithTimeout method.
func GoWaitResult(fun CancelCallable) FutureTask {
	task := WrapWaitResultTask(fun)
	go task.Run()
	return task
}

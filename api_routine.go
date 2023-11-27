package routine

// Runnable provides a function without return values.
type Runnable func()

// Callable provides a function that returns a value of type TResult.
type Callable[TResult any] func() TResult

// CancelRunnable provides a cancellable function without return values.
type CancelRunnable func(token CancelToken)

// CancelCallable provides a cancellable function that returns a value of type TResult.
type CancelCallable[TResult any] func(token CancelToken) TResult

// WrapTask create a new task and capture the inheritableThreadLocals from the current goroutine.
// This function returns a FutureTask instance, but the return task will not run automatically.
// You can run it in a sub-goroutine or goroutine-pool by FutureTask.Run method, wait by FutureTask.Get or FutureTask.GetWithTimeout method.
// When the returned task run panic will be caught and error stack will be printed, the panic will be trigger again when calling FutureTask.Get or FutureTask.GetWithTimeout method.
func WrapTask(fun Runnable) FutureTask[any] {
	ctx := createInheritedMap()
	callable := inheritedTask{context: ctx, function: fun}.run
	return NewFutureTask[any](callable)
}

// WrapWaitTask create a new task and capture the inheritableThreadLocals from the current goroutine.
// This function returns a FutureTask instance, but the return task will not run automatically.
// You can run it in a sub-goroutine or goroutine-pool by FutureTask.Run method, wait by FutureTask.Get or FutureTask.GetWithTimeout method.
// When the returned task run panic will be caught, the panic will be trigger again when calling FutureTask.Get or FutureTask.GetWithTimeout method.
func WrapWaitTask(fun CancelRunnable) FutureTask[any] {
	ctx := createInheritedMap()
	callable := inheritedWaitTask{context: ctx, function: fun}.run
	return NewFutureTask[any](callable)
}

// WrapWaitResultTask create a new task and capture the inheritableThreadLocals from the current goroutine.
// This function returns a FutureTask instance, but the return task will not run automatically.
// You can run it in a sub-goroutine or goroutine-pool by FutureTask.Run method, wait and get result by FutureTask.Get or FutureTask.GetWithTimeout method.
// When the returned task run panic will be caught, the panic will be trigger again when calling FutureTask.Get or FutureTask.GetWithTimeout method.
func WrapWaitResultTask[TResult any](fun CancelCallable[TResult]) FutureTask[TResult] {
	ctx := createInheritedMap()
	callable := inheritedWaitResultTask[TResult]{context: ctx, function: fun}.run
	return NewFutureTask[TResult](callable)
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
func GoWait(fun CancelRunnable) FutureTask[any] {
	task := WrapWaitTask(fun)
	go task.Run()
	return task
}

// GoWaitResult starts a new goroutine, and copy inheritableThreadLocals from current goroutine.
// This function will auto invoke the func and return a FutureTask instance, so we can wait and get result by FutureTask.Get or FutureTask.GetWithTimeout method.
// If panic occur in goroutine, The panic will be trigger again when calling FutureTask.Get or FutureTask.GetWithTimeout method.
func GoWaitResult[TResult any](fun CancelCallable[TResult]) FutureTask[TResult] {
	task := WrapWaitResultTask(fun)
	go task.Run()
	return task
}

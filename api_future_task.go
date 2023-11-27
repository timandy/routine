package routine

import "time"

// FutureCallable provides a future function that returns a value of type TResult.
type FutureCallable[TResult any] func(task FutureTask[TResult]) TResult

// CancelToken propagates notification that operations should be canceled.
type CancelToken interface {
	// IsCanceled returns true if task was canceled.
	IsCanceled() bool

	// Cancel notifies the waiting coroutine that the task has canceled and returns stack information.
	Cancel()
}

// FutureTask provide a way to wait for the sub-coroutine to finish executing, get the return value of the sub-coroutine, and catch the sub-coroutine panic.
type FutureTask[TResult any] interface {
	// IsDone returns true if completed in any fashion: normally, exceptionally or via cancellation.
	IsDone() bool

	// IsCanceled returns true if task was canceled.
	IsCanceled() bool

	// IsFailed returns true if completed exceptionally.
	IsFailed() bool

	// Complete notifies the waiting coroutine that the task has completed normally and returns the execution result.
	Complete(result TResult)

	// Cancel notifies the waiting coroutine that the task has canceled and returns stack information.
	Cancel()

	// Fail notifies the waiting coroutine that the task has terminated due to panic and returns stack information.
	Fail(error any)

	// Get return the execution result of the sub-coroutine, if there is no result, return nil.
	// If task is canceled, a panic with cancellation will be raised.
	// If panic is raised during the execution of the sub-coroutine, it will be raised again at this time.
	Get() TResult

	// GetWithTimeout return the execution result of the sub-coroutine, if there is no result, return nil.
	// If task is canceled, a panic with cancellation will be raised.
	// If panic is raised during the execution of the sub-coroutine, it will be raised again at this time.
	// If the deadline is reached, a panic with timeout error will be raised.
	GetWithTimeout(timeout time.Duration) TResult

	// Run execute the task, the method can be called repeatedly, but the task will only execute once.
	Run()
}

// NewFutureTask Create a new instance.
func NewFutureTask[TResult any](callable FutureCallable[TResult]) FutureTask[TResult] {
	if callable == nil {
		panic("callable can not be nil.")
	}
	task := &futureTask[TResult]{callable: callable}
	task.await.Add(1)
	return task
}

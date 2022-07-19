package routine

import "time"

// Future provide a way to wait for the sub-coroutine to finish executing, get the return value of the sub-coroutine, and catch the sub-coroutine panic.
type Future interface {
	// IsDone returns true if completed in any fashion: normally, exceptionally or via cancellation.
	IsDone() bool

	// IsCanceled returns true if task be canceled.
	IsCanceled() bool

	// IsFailed returns true if completed exceptionally.
	IsFailed() bool

	// Complete notifies the parent coroutine that the task has completed normally and returns the execution result.
	Complete(result any)

	// Cancel notifies the parent coroutine that the task is canceled due to special reason and returns stack information.
	Cancel(reason any)

	// Fail notifies the parent coroutine that the task is terminated due to panic and returns stack information.
	Fail(error any)

	// Get return the execution result of the sub-coroutine, if there is no result, return nil.
	// If task is canceled, a panic with cancel reason will be raised.
	// If panic is raised during the execution of the sub-coroutine, it will be raised again at this time.
	Get() any

	// GetWithTimeout return the execution result of the sub-coroutine, if there is no result, return nil.
	// If task is canceled, a panic with cancel reason will be raised.
	// If panic is raised during the execution of the sub-coroutine, it will be raised again at this time.
	// If the deadline is reached, a panic with timeout error will be raised.
	GetWithTimeout(timeout time.Duration) any
}

// NewFuture Create a new instance.
func NewFuture() Future {
	fut := &future{}
	fut.await.Add(1)
	return fut
}

package routine

import "sync"

// Feature provide a way to wait for the sub-coroutine to finish executing, get the return value of the sub-coroutine, and catch the sub-coroutine panic.
type Feature interface {
	// Complete notifies the parent coroutine that the task has completed and returns the execution result.
	// This method is called by the child coroutine.
	Complete(result Any)

	// CompleteError notifies the parent coroutine that the task is terminated due to panic and returns stack information.
	// This method is called by the child coroutine.
	CompleteError(error Any)

	// Get the execution result of the sub-coroutine, if there is no result, return nil.
	// If panic is raised during the execution of the sub-coroutine, it will be raised again at this time.
	// this method is called by the parent coroutine.
	Get() Any
}

// NewFeature Create a new instance.
func NewFeature() Feature {
	await := &sync.WaitGroup{}
	await.Add(1)
	return &feature{await: await}
}

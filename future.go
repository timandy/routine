package routine

import "sync"

type future struct {
	await  *sync.WaitGroup
	error  RuntimeError
	result Any
}

func (fut *future) Complete(result Any) {
	fut.result = result
	fut.await.Done()
}

func (fut *future) CompleteError(error Any) {
	if runtimeErr, isRuntimeErr := error.(RuntimeError); isRuntimeErr {
		fut.error = runtimeErr
	} else {
		fut.error = NewRuntimeError(error)
	}
	fut.await.Done()
}

func (fut *future) Get() Any {
	fut.await.Wait()
	if fut.error != nil {
		panic(fut.error)
	}
	return fut.result
}

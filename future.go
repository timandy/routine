package routine

import "sync"

type future struct {
	await  *sync.WaitGroup
	error  RuntimeError
	result Any
}

func (fea *future) Complete(result Any) {
	fea.result = result
	fea.await.Done()
}

func (fea *future) CompleteError(error Any) {
	if runtimeErr, isRuntimeErr := error.(RuntimeError); isRuntimeErr {
		fea.error = runtimeErr
	} else {
		fea.error = NewRuntimeError(error)
	}
	fea.await.Done()
}

func (fea *future) Get() Any {
	fea.await.Wait()
	if fea.error != nil {
		panic(fea.error)
	}
	return fea.result
}

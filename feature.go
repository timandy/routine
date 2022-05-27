package routine

import "sync"

type feature struct {
	await  *sync.WaitGroup
	error  RuntimeError
	result Any
}

func (fea *feature) Complete(result Any) {
	fea.result = result
	fea.await.Done()
}

func (fea *feature) CompleteError(error Any) {
	if runtimeErr, isRuntimeErr := error.(RuntimeError); isRuntimeErr {
		fea.error = runtimeErr
	} else {
		fea.error = NewRuntimeError(error)
	}
	fea.await.Done()
}

func (fea *feature) Get() Any {
	fea.await.Wait()
	if fea.error != nil {
		panic(fea.error)
	}
	return fea.result
}

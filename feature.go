package routine

import "sync"

type feature struct {
	await  *sync.WaitGroup
	error  StackError
	result Any
}

func (fea *feature) Complete(result Any) {
	fea.result = result
	fea.await.Done()
}

func (fea *feature) CompleteError(error Any) {
	fea.error = NewStackError(error)
	fea.await.Done()
}

func (fea *feature) Get() Any {
	fea.await.Wait()
	if fea.error != nil {
		panic(fea.error)
	}
	return fea.result
}

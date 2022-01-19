package routine

import "sync"

type feature struct {
	waitGroup *sync.WaitGroup
	error     *StackError
	result    Any
}

func (fea *feature) Complete(result Any) {
	fea.result = result
	fea.waitGroup.Done()
}

func (fea *feature) CompleteError(error Any) {
	fea.error = NewStackError(error)
	fea.waitGroup.Done()
}

func (fea *feature) Get() Any {
	fea.waitGroup.Wait()
	if fea.error != nil {
		panic(fea.error)
	}
	return fea.result
}

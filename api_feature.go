package routine

import "sync"

type Feature struct {
	waitGroup *sync.WaitGroup
	error     *StackError
	result    Any
}

func NewFeature() *Feature {
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)
	return &Feature{waitGroup: waitGroup}
}

func (fea *Feature) Complete(result Any) {
	fea.result = result
	fea.waitGroup.Done()
}

func (fea *Feature) CompleteError(error Any) {
	fea.error = NewStackError(error)
	fea.waitGroup.Done()
}

func (fea *Feature) Get() Any {
	fea.waitGroup.Wait()
	if fea.error != nil {
		panic(fea.error)
	}
	return fea.result
}

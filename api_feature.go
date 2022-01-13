package routine

import "sync"

type Feature struct {
	waitGroup *sync.WaitGroup
	error     *StackError
	result    Any
}

func feature() *Feature {
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)
	return &Feature{waitGroup: waitGroup}
}

func (f *Feature) complete(result Any) {
	f.result = result
	f.waitGroup.Done()
}

func (f *Feature) completeError(error *StackError) {
	f.error = error
	f.waitGroup.Done()
}

func (f *Feature) Get() Any {
	f.waitGroup.Wait()
	if f.error != nil {
		panic(f.error)
	}
	return f.result
}

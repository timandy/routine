package routine

import "sync"

type Feature interface {
	Complete(result Any)

	CompleteError(error Any)

	Get() Any
}

func NewFeature() Feature {
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)
	return &feature{waitGroup: waitGroup}
}

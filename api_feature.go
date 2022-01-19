package routine

import "sync"

type Feature interface {
	Complete(result Any)

	CompleteError(error Any)

	Get() Any
}

func NewFeature() Feature {
	await := &sync.WaitGroup{}
	await.Add(1)
	return &feature{await: await}
}

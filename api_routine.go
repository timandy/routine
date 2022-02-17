package routine

import "fmt"

// Go starts a new goroutine.
// This function will auto invoke the fun and print error stack when panic occur in goroutine.
func Go(fun func()) {
	go func() {
		// catch
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(NewStackError(err).Error())
			}
		}()
		// invoke
		fun()
	}()
}

// GoWait starts a new goroutine.
// This function return a Feature pointer, so we can wait by Feature.Get method.
// If panic occur in goroutine, The panic will be trigger again when calling Feature.Get method.
func GoWait(fun func()) Feature {
	fea := NewFeature()
	go func() {
		// catch
		defer func() {
			if err := recover(); err != nil {
				fea.CompleteError(err)
			}
		}()
		// invoke
		fun()
		fea.Complete(nil)
	}()
	return fea
}

// GoWaitResult starts a new goroutine.
// This function return a Feature pointer, so we can wait and get result by Feature.Get method.
// If panic occur in goroutine, The panic will be trigger again when calling Feature.Get method.
func GoWaitResult(fun func() Any) Feature {
	fea := NewFeature()
	go func() {
		// catch
		defer func() {
			if err := recover(); err != nil {
				fea.CompleteError(err)
			}
		}()
		// invoke
		fea.Complete(fun())
	}()
	return fea
}

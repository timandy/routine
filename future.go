package routine

import (
	"fmt"
	"sync"
	"time"
)

type futureStatus int

const (
	running futureStatus = iota
	completed
	canceled
	failed
)

type future[T any] struct {
	lock   sync.RWMutex
	await  sync.WaitGroup
	status futureStatus
	error  RuntimeError
	result T
}

func (fut *future[T]) IsDone() bool {
	fut.lock.RLock()
	defer fut.lock.RUnlock()
	return fut.status != running
}

func (fut *future[T]) IsCanceled() bool {
	fut.lock.RLock()
	defer fut.lock.RUnlock()
	return fut.status == canceled
}

func (fut *future[T]) IsFailed() bool {
	fut.lock.RLock()
	defer fut.lock.RUnlock()
	return fut.status == failed
}

func (fut *future[T]) Complete(result T) {
	fut.lock.Lock()
	defer fut.lock.Unlock()
	if fut.status != running {
		return
	}
	fut.result = result
	fut.status = completed
	fut.await.Done()
}

func (fut *future[T]) Cancel() {
	fut.lock.Lock()
	defer fut.lock.Unlock()
	if fut.status != running {
		return
	}
	fut.error = NewRuntimeError("Task was canceled.")
	fut.status = canceled
	fut.await.Done()
}

func (fut *future[T]) Fail(error any) {
	fut.lock.Lock()
	defer fut.lock.Unlock()
	if fut.status != running {
		return
	}
	runtimeErr, isRuntimeErr := error.(RuntimeError)
	if !isRuntimeErr {
		runtimeErr = NewRuntimeError(error)
	}
	fut.error = runtimeErr
	fut.status = failed
	fut.await.Done()
}

func (fut *future[T]) Get() T {
	fut.await.Wait()
	if fut.status == completed {
		return fut.result
	}
	panic(fut.error)
}

func (fut *future[T]) GetWithTimeout(timeout time.Duration) T {
	resultChan := make(chan struct{})
	errorChan := make(chan struct{})
	go func() {
		fut.await.Wait()
		if fut.status == completed {
			close(resultChan)
			return
		}
		close(errorChan)
	}()
	select {
	case <-resultChan:
		return fut.result
	case <-errorChan:
		panic(fut.error)
	case <-time.After(timeout):
		fut.timeout(timeout)
		fut.await.Wait()
		if fut.status == completed {
			return fut.result
		}
		panic(fut.error)
	}
}

func (fut *future[T]) timeout(timeout time.Duration) {
	fut.lock.Lock()
	defer fut.lock.Unlock()
	if fut.status != running {
		return
	}
	fut.error = NewRuntimeError(fmt.Sprintf("Task execution timeout after %v.", timeout))
	fut.status = canceled
	fut.await.Done()
}

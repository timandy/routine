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

type future struct {
	lock   sync.RWMutex
	await  sync.WaitGroup
	status futureStatus
	error  RuntimeError
	result any
}

func (fut *future) IsDone() bool {
	fut.lock.RLock()
	defer fut.lock.RUnlock()
	return fut.status != running
}

func (fut *future) IsCanceled() bool {
	fut.lock.RLock()
	defer fut.lock.RUnlock()
	return fut.status == canceled
}

func (fut *future) IsFailed() bool {
	fut.lock.RLock()
	defer fut.lock.RUnlock()
	return fut.status == failed
}

func (fut *future) Complete(result any) {
	fut.lock.Lock()
	defer fut.lock.Unlock()
	if fut.status != running {
		return
	}
	fut.result = result
	fut.status = completed
	fut.await.Done()
}

func (fut *future) Cancel(reason any) {
	fut.lock.Lock()
	defer fut.lock.Unlock()
	if fut.status != running {
		return
	}
	runtimeErr, isRuntimeErr := reason.(RuntimeError)
	if !isRuntimeErr {
		runtimeErr = NewRuntimeError(reason)
	}
	fut.error = runtimeErr
	fut.status = canceled
	fut.await.Done()
}

func (fut *future) Fail(error any) {
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

func (fut *future) Get() any {
	fut.await.Wait()
	if fut.status == completed {
		return fut.result
	}
	panic(fut.error)
}

func (fut *future) GetWithTimeout(timeout time.Duration) any {
	resultChan := make(chan struct{}, 1)
	errorChan := make(chan struct{}, 1)
	go func() {
		defer func() {
			close(resultChan)
			close(errorChan)
		}()
		fut.await.Wait()
		if fut.status == completed {
			resultChan <- struct{}{}
			return
		}
		errorChan <- struct{}{}
	}()
	select {
	case <-resultChan:
		return fut.result
	case <-errorChan:
		panic(fut.error)
	case <-time.After(timeout):
		timeoutError := NewRuntimeErrorWithMessage(fmt.Sprintf("Task execution timeout after %v.", timeout))
		fut.Cancel(timeoutError)
		panic(timeoutError)
	}
}

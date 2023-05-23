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

type futureTask struct {
	lock   sync.RWMutex
	await  sync.WaitGroup
	status futureStatus
	error  RuntimeError
	result any
}

func (task *futureTask) IsDone() bool {
	task.lock.RLock()
	defer task.lock.RUnlock()
	return task.status != running
}

func (task *futureTask) IsCanceled() bool {
	task.lock.RLock()
	defer task.lock.RUnlock()
	return task.status == canceled
}

func (task *futureTask) IsFailed() bool {
	task.lock.RLock()
	defer task.lock.RUnlock()
	return task.status == failed
}

func (task *futureTask) Complete(result any) {
	task.lock.Lock()
	defer task.lock.Unlock()
	if task.status != running {
		return
	}
	task.result = result
	task.status = completed
	task.await.Done()
}

func (task *futureTask) Cancel() {
	task.lock.Lock()
	defer task.lock.Unlock()
	if task.status != running {
		return
	}
	task.error = NewRuntimeError("Task was canceled.")
	task.status = canceled
	task.await.Done()
}

func (task *futureTask) Fail(error any) {
	task.lock.Lock()
	defer task.lock.Unlock()
	if task.status != running {
		return
	}
	runtimeErr, isRuntimeErr := error.(RuntimeError)
	if !isRuntimeErr {
		runtimeErr = NewRuntimeError(error)
	}
	task.error = runtimeErr
	task.status = failed
	task.await.Done()
}

func (task *futureTask) Get() any {
	task.await.Wait()
	if task.status == completed {
		return task.result
	}
	panic(task.error)
}

func (task *futureTask) GetWithTimeout(timeout time.Duration) any {
	waitChan := make(chan struct{})
	go func() {
		task.await.Wait()
		close(waitChan)
	}()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-waitChan:
		if task.status == completed {
			return task.result
		}
		panic(task.error)
	case <-timer.C:
		task.timeout(timeout)
		task.await.Wait()
		if task.status == completed {
			return task.result
		}
		panic(task.error)
	}
}

func (task *futureTask) timeout(timeout time.Duration) {
	task.lock.Lock()
	defer task.lock.Unlock()
	if task.status != running {
		return
	}
	task.error = NewRuntimeError(fmt.Sprintf("Task execution timeout after %v.", timeout))
	task.status = canceled
	task.await.Done()
}

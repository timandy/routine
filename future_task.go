package routine

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type taskState = int32

const (
	taskStateNew taskState = iota
	taskStateRunning
	taskStateCompleted
	taskStateCanceled
	taskStateFailed
)

type futureTask[TResult any] struct {
	await    sync.WaitGroup
	state    taskState
	callable FutureCallable[TResult]
	result   TResult
	error    RuntimeError
}

func (task *futureTask[TResult]) IsDone() bool {
	state := atomic.LoadInt32(&task.state)
	return state == taskStateCompleted || state == taskStateCanceled || state == taskStateFailed
}

func (task *futureTask[TResult]) IsCanceled() bool {
	return atomic.LoadInt32(&task.state) == taskStateCanceled
}

func (task *futureTask[TResult]) IsFailed() bool {
	return atomic.LoadInt32(&task.state) == taskStateFailed
}

func (task *futureTask[TResult]) Complete(result TResult) {
	if atomic.CompareAndSwapInt32(&task.state, taskStateNew, taskStateCompleted) ||
		atomic.CompareAndSwapInt32(&task.state, taskStateRunning, taskStateCompleted) {
		task.result = result
		task.await.Done()
	}
}

func (task *futureTask[TResult]) Cancel() {
	if atomic.CompareAndSwapInt32(&task.state, taskStateNew, taskStateCanceled) ||
		atomic.CompareAndSwapInt32(&task.state, taskStateRunning, taskStateCanceled) {
		task.error = NewRuntimeError("Task was canceled.")
		task.await.Done()
	}
}

func (task *futureTask[TResult]) Fail(error any) {
	if atomic.CompareAndSwapInt32(&task.state, taskStateNew, taskStateFailed) ||
		atomic.CompareAndSwapInt32(&task.state, taskStateRunning, taskStateFailed) {
		runtimeErr, isRuntimeErr := error.(RuntimeError)
		if !isRuntimeErr {
			runtimeErr = NewRuntimeError(error)
		}
		task.error = runtimeErr
		task.await.Done()
	}
}

func (task *futureTask[TResult]) Get() TResult {
	task.await.Wait()
	if atomic.LoadInt32(&task.state) == taskStateCompleted {
		return task.result
	}
	panic(task.error)
}

func (task *futureTask[TResult]) GetWithTimeout(timeout time.Duration) TResult {
	waitChan := make(chan struct{})
	go func() {
		task.await.Wait()
		close(waitChan)
	}()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-waitChan:
		if atomic.LoadInt32(&task.state) == taskStateCompleted {
			return task.result
		}
		panic(task.error)
	case <-timer.C:
		task.timeout(timeout)
		task.await.Wait()
		if atomic.LoadInt32(&task.state) == taskStateCompleted {
			return task.result
		}
		panic(task.error)
	}
}

func (task *futureTask[TResult]) Run() {
	if atomic.CompareAndSwapInt32(&task.state, taskStateNew, taskStateRunning) {
		defer func() {
			if cause := recover(); cause != nil {
				task.Fail(cause)
			}
		}()
		result := task.callable(task)
		task.Complete(result)
	}
}

func (task *futureTask[TResult]) timeout(timeout time.Duration) {
	if atomic.CompareAndSwapInt32(&task.state, taskStateNew, taskStateCanceled) ||
		atomic.CompareAndSwapInt32(&task.state, taskStateRunning, taskStateCanceled) {
		task.error = NewRuntimeError(fmt.Sprintf("Task execution timeout after %v.", timeout))
		task.await.Done()
	}
}

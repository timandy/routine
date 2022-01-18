package routine

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestGo_Error(t *testing.T) {
	run := false
	assert.NotPanics(t, func() {
		waiter := &sync.WaitGroup{}
		waiter.Add(1)
		Go(func() {
			run = true
			waiter.Done()
			panic("error")
		})
		waiter.Wait()
	})
	assert.True(t, run)
}

func TestGo_Nil(t *testing.T) {
	copied := createInheritedMap()
	assert.Nil(t, copied)
	//
	run := false
	waiter := &sync.WaitGroup{}
	waiter.Add(1)
	Go(func() {
		thd := currentThread(true)
		assert.Nil(t, thd.inheritableThreadLocals)
		run = true
		waiter.Done()
	})
	waiter.Wait()
	assert.True(t, run)
}

func TestGo_Value(t *testing.T) {
	threadLocal := NewThreadLocal()
	threadLocal.Set("Hello")
	assert.Equal(t, "Hello", threadLocal.Get())
	//
	inheritableThreadLocal := NewInheritableThreadLocal()
	inheritableThreadLocal.Set("World")
	assert.Equal(t, "World", inheritableThreadLocal.Get())
	//
	copied := createInheritedMap()
	assert.NotNil(t, copied)
	//
	run := false
	waiter := &sync.WaitGroup{}
	waiter.Add(1)
	Go(func() {
		thd := currentThread(true)
		assert.NotNil(t, thd.inheritableThreadLocals == nil)
		//
		assert.Nil(t, threadLocal.Get())
		assert.Equal(t, "World", inheritableThreadLocal.Get())
		//
		threadLocal.Set("Hello2")
		assert.Equal(t, "Hello2", threadLocal.Get())
		//
		inheritableThreadLocal.Remove()
		assert.Nil(t, inheritableThreadLocal.Get())
		//
		run = true
		waiter.Done()
	})
	waiter.Wait()
	assert.True(t, run)
	//
	assert.Equal(t, "Hello", threadLocal.Get())
	assert.Equal(t, "World", inheritableThreadLocal.Get())
}

func TestGoWait_Error(t *testing.T) {
	run := false
	assert.Panics(t, func() {
		fea := GoWait(func() {
			run = true
			panic("error")
		})
		fea.Get()
	})
	assert.True(t, run)
}

func TestGoWait_Nil(t *testing.T) {
	copied := createInheritedMap()
	assert.Nil(t, copied)
	//
	run := false
	fea := GoWait(func() {
		thd := currentThread(true)
		assert.Nil(t, thd.inheritableThreadLocals)
		run = true
	})
	assert.Nil(t, fea.Get())
	assert.True(t, run)
}

func TestGoWait_Value(t *testing.T) {
	threadLocal := NewThreadLocal()
	threadLocal.Set("Hello")
	assert.Equal(t, "Hello", threadLocal.Get())
	//
	inheritableThreadLocal := NewInheritableThreadLocal()
	inheritableThreadLocal.Set("World")
	assert.Equal(t, "World", inheritableThreadLocal.Get())
	//
	copied := createInheritedMap()
	assert.NotNil(t, copied)
	//
	run := false
	fea := GoWait(func() {
		thd := currentThread(true)
		assert.NotNil(t, thd.inheritableThreadLocals == nil)
		//
		assert.Nil(t, threadLocal.Get())
		assert.Equal(t, "World", inheritableThreadLocal.Get())
		//
		threadLocal.Set("Hello2")
		assert.Equal(t, "Hello2", threadLocal.Get())
		//
		inheritableThreadLocal.Remove()
		assert.Nil(t, inheritableThreadLocal.Get())
		//
		run = true
	})
	assert.Nil(t, fea.Get())
	assert.True(t, run)
	//
	assert.Equal(t, "Hello", threadLocal.Get())
	assert.Equal(t, "World", inheritableThreadLocal.Get())
}

func TestGoWaitResult_Error(t *testing.T) {
	run := false
	assert.Panics(t, func() {
		fea := GoWaitResult(func() Any {
			run = true
			if run {
				panic("error")
			}
			return 1
		})
		fea.Get()
	})
	assert.True(t, run)
}

func TestGoWaitResult_Nil(t *testing.T) {
	copied := createInheritedMap()
	assert.Nil(t, copied)
	//
	run := false
	fea := GoWaitResult(func() Any {
		thd := currentThread(true)
		assert.Nil(t, thd.inheritableThreadLocals)
		run = true
		return true
	})
	assert.True(t, fea.Get().(bool))
	assert.True(t, run)
}

func TestGoWaitResult_Value(t *testing.T) {
	threadLocal := NewThreadLocal()
	threadLocal.Set("Hello")
	assert.Equal(t, "Hello", threadLocal.Get())
	//
	inheritableThreadLocal := NewInheritableThreadLocal()
	inheritableThreadLocal.Set("World")
	assert.Equal(t, "World", inheritableThreadLocal.Get())
	//
	copied := createInheritedMap()
	assert.NotNil(t, copied)
	//
	run := false
	fea := GoWaitResult(func() Any {
		thd := currentThread(true)
		assert.NotNil(t, thd.inheritableThreadLocals == nil)
		//
		assert.Nil(t, threadLocal.Get())
		assert.Equal(t, "World", inheritableThreadLocal.Get())
		//
		threadLocal.Set("Hello2")
		assert.Equal(t, "Hello2", threadLocal.Get())
		//
		inheritableThreadLocal.Remove()
		assert.Nil(t, inheritableThreadLocal.Get())
		//
		run = true
		return true
	})
	assert.True(t, fea.Get().(bool))
	assert.True(t, run)
	//
	assert.Equal(t, "Hello", threadLocal.Get())
	assert.Equal(t, "World", inheritableThreadLocal.Get())
}

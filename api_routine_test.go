package routine

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestGo_Error(t *testing.T) {
	run := false
	assert.NotPanics(t, func() {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		Go(func() {
			run = true
			wg.Done()
			panic("error")
		})
		wg.Wait()
	})
	assert.True(t, run)
}

func TestGo_Nil(t *testing.T) {
	assert.Nil(t, createInheritedMap())
	//
	run := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	Go(func() {
		assert.Nil(t, createInheritedMap())
		run = true
		wg.Done()
	})
	wg.Wait()
	assert.True(t, run)
}

func TestGo_Value(t *testing.T) {
	tls := NewThreadLocal()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	inheritableTls := NewInheritableThreadLocal()
	inheritableTls.Set("World")
	assert.Equal(t, "World", inheritableTls.Get())
	//
	assert.NotNil(t, createInheritedMap())
	//
	run := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	Go(func() {
		assert.NotNil(t, createInheritedMap())
		//
		assert.Nil(t, tls.Get())
		assert.Equal(t, "World", inheritableTls.Get())
		//
		tls.Set("Hello2")
		assert.Equal(t, "Hello2", tls.Get())
		//
		inheritableTls.Remove()
		assert.Nil(t, inheritableTls.Get())
		//
		run = true
		wg.Done()
	})
	wg.Wait()
	assert.True(t, run)
	//
	assert.Equal(t, "Hello", tls.Get())
	assert.Equal(t, "World", inheritableTls.Get())
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
	assert.Nil(t, createInheritedMap())
	//
	run := false
	fea := GoWait(func() {
		assert.Nil(t, createInheritedMap())
		run = true
	})
	assert.Nil(t, fea.Get())
	assert.True(t, run)
}

func TestGoWait_Value(t *testing.T) {
	tls := NewThreadLocal()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	inheritableTls := NewInheritableThreadLocal()
	inheritableTls.Set("World")
	assert.Equal(t, "World", inheritableTls.Get())
	//
	assert.NotNil(t, createInheritedMap())
	//
	run := false
	fea := GoWait(func() {
		assert.NotNil(t, createInheritedMap())
		//
		assert.Nil(t, tls.Get())
		assert.Equal(t, "World", inheritableTls.Get())
		//
		tls.Set("Hello2")
		assert.Equal(t, "Hello2", tls.Get())
		//
		inheritableTls.Remove()
		assert.Nil(t, inheritableTls.Get())
		//
		run = true
	})
	assert.Nil(t, fea.Get())
	assert.True(t, run)
	//
	assert.Equal(t, "Hello", tls.Get())
	assert.Equal(t, "World", inheritableTls.Get())
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
	assert.Nil(t, createInheritedMap())
	//
	run := false
	fea := GoWaitResult(func() Any {
		assert.Nil(t, createInheritedMap())
		run = true
		return true
	})
	assert.True(t, fea.Get().(bool))
	assert.True(t, run)
}

func TestGoWaitResult_Value(t *testing.T) {
	tls := NewThreadLocal()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	inheritableTls := NewInheritableThreadLocal()
	inheritableTls.Set("World")
	assert.Equal(t, "World", inheritableTls.Get())
	//
	assert.NotNil(t, createInheritedMap())
	//
	run := false
	fea := GoWaitResult(func() Any {
		assert.NotNil(t, createInheritedMap())
		//
		assert.Nil(t, tls.Get())
		assert.Equal(t, "World", inheritableTls.Get())
		//
		tls.Set("Hello2")
		assert.Equal(t, "Hello2", tls.Get())
		//
		inheritableTls.Remove()
		assert.Nil(t, inheritableTls.Get())
		//
		run = true
		return true
	})
	assert.True(t, fea.Get().(bool))
	assert.True(t, run)
	//
	assert.Equal(t, "Hello", tls.Get())
	assert.Equal(t, "World", inheritableTls.Get())
}

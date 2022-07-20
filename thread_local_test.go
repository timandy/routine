package routine

import (
	"math"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThreadLocal_Index(t *testing.T) {
	tls := NewThreadLocal()
	assert.GreaterOrEqual(t, tls.(*threadLocal).index, 0)
	tls2 := NewThreadLocalWithInitial(func() any {
		return "Hello"
	})
	assert.Greater(t, tls2.(*threadLocal).index, tls.(*threadLocal).index)
}

func TestThreadLocal_NextIndex(t *testing.T) {
	backup := threadLocalIndex
	defer func() {
		threadLocalIndex = backup
	}()
	//
	threadLocalIndex = math.MaxInt32
	assert.Panics(t, func() {
		nextThreadLocalIndex()
	})
	assert.Equal(t, math.MaxInt32, int(threadLocalIndex))
}

func TestThreadLocal_Common(t *testing.T) {
	tls := NewThreadLocal()
	tls2 := NewThreadLocal()
	tls.Remove()
	tls2.Remove()
	assert.Nil(t, tls.Get())
	assert.Nil(t, tls2.Get())
	//
	tls.Set(1)
	tls2.Set("World")
	assert.Equal(t, 1, tls.Get())
	assert.Equal(t, "World", tls2.Get())
	//
	tls.Set(nil)
	tls2.Set(nil)
	assert.Nil(t, tls.Get())
	assert.Nil(t, tls2.Get())
	//
	tls.Set(2)
	tls2.Set("!")
	assert.Equal(t, 2, tls.Get())
	assert.Equal(t, "!", tls2.Get())
	//
	tls.Remove()
	tls2.Remove()
	assert.Nil(t, tls.Get())
	assert.Nil(t, tls2.Get())
	//
	tls.Set(2)
	tls2.Set("!")
	assert.Equal(t, 2, tls.Get())
	assert.Equal(t, "!", tls2.Get())
	wg := &sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		Go(func() {
			assert.Nil(t, tls.Get())
			assert.Nil(t, tls2.Get())
			wg.Done()
		})
	}
	wg.Wait()
	assert.Equal(t, 2, tls.Get())
	assert.Equal(t, "!", tls2.Get())
}

func TestThreadLocal_Mixed(t *testing.T) {
	tls := NewThreadLocal()
	tls2 := NewThreadLocalWithInitial(func() any {
		return "Hello"
	})
	assert.Nil(t, tls.Get())
	assert.Equal(t, "Hello", tls2.Get())
	//
	tls.Set(1)
	tls2.Set("World")
	assert.Equal(t, 1, tls.Get())
	assert.Equal(t, "World", tls2.Get())
	//
	tls.Set(nil)
	tls2.Set(nil)
	assert.Nil(t, tls.Get())
	assert.Nil(t, tls2.Get())
	//
	tls.Set(2)
	tls2.Set("!")
	assert.Equal(t, 2, tls.Get())
	assert.Equal(t, "!", tls2.Get())
	//
	tls.Remove()
	tls2.Remove()
	assert.Nil(t, tls.Get())
	assert.Equal(t, "Hello", tls2.Get())
	//
	tls.Set(2)
	tls2.Set("!")
	assert.Equal(t, 2, tls.Get())
	assert.Equal(t, "!", tls2.Get())
	wg := &sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		Go(func() {
			assert.Nil(t, tls.Get())
			assert.Equal(t, "Hello", tls2.Get())
			wg.Done()
		})
	}
	wg.Wait()
	assert.Equal(t, 2, tls.Get())
	assert.Equal(t, "!", tls2.Get())
}

func TestThreadLocal_WithInitial(t *testing.T) {
	src := &person{Id: 1, Name: "Tim"}
	tls := NewThreadLocalWithInitial(nil)
	tls2 := NewThreadLocalWithInitial(func() any {
		return nil
	})
	tls3 := NewThreadLocalWithInitial(func() any {
		return src
	})
	tls4 := NewThreadLocalWithInitial(func() any {
		return *src
	})

	for i := 0; i < 100; i++ {
		p := tls.Get()
		assert.Nil(t, p)
		//
		p2 := tls2.Get()
		assert.Nil(t, p2)
		//
		p3 := tls3.Get().(*person)
		assert.Same(t, src, p3)

		p4 := tls4.Get().(person)
		assert.NotSame(t, src, &p4)
		assert.Equal(t, *src, p4)

		wg := &sync.WaitGroup{}
		wg.Add(1)
		Go(func() {
			assert.Same(t, src, tls3.Get().(*person))
			p5 := tls4.Get().(person)
			assert.NotSame(t, src, &p5)
			assert.Equal(t, *src, p5)
			//
			wg.Done()
		})
		wg.Wait()
	}

	tls3.Set(nil)
	tls4.Set(nil)
	assert.Nil(t, tls3.Get())
	assert.Nil(t, tls4.Get())

	tls3.Remove()
	tls4.Remove()
	assert.Same(t, src, tls3.Get().(*person))
	p6 := tls4.Get().(person)
	assert.NotSame(t, src, &p6)
	assert.Equal(t, *src, p6)
}

func TestThreadLocal_CrossCoroutine(t *testing.T) {
	tls := NewThreadLocal()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get().(string))
	subWait := &sync.WaitGroup{}
	subWait.Add(2)
	finishWait := &sync.WaitGroup{}
	finishWait.Add(2)
	go func() {
		subWait.Wait()
		assert.Nil(t, tls.Get())
		finishWait.Done()
	}()
	Go(func() {
		subWait.Wait()
		assert.Nil(t, tls.Get())
		finishWait.Done()
	})
	tls.Remove()      //remove in parent goroutine should not affect child goroutine
	subWait.Done()    //allow sub goroutine run
	subWait.Done()    //allow sub goroutine run
	finishWait.Wait() //wait sub goroutine done
	finishWait.Wait() //wait sub goroutine done
}

func TestThreadLocal_CreateBatch(t *testing.T) {
	const count = 128
	tlsList := make([]ThreadLocal, count)
	for i := 0; i < count; i++ {
		value := i
		tlsList[i] = NewThreadLocalWithInitial(func() any { return value })
	}
	for i := 0; i < count; i++ {
		assert.Equal(t, i, tlsList[i].Get())
	}
}

type person struct {
	Id   int
	Name string
}

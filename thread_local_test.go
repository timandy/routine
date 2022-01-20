package routine

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestThreadLocal_Id(t *testing.T) {
	tls := NewThreadLocal()
	assert.GreaterOrEqual(t, tls.Id(), 0)
	tls2 := NewThreadLocalWithInitial(func() Any {
		return "Hello"
	})
	assert.Greater(t, tls2.Id(), tls.Id())
}

func TestThreadLocal(t *testing.T) {
	tls := NewThreadLocal()
	tls2 := NewThreadLocal()
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

func TestThreadLocalMixed(t *testing.T) {
	tls := NewThreadLocal()
	tls2 := NewThreadLocalWithInitial(func() Any {
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

func TestThreadLocalWithInitial(t *testing.T) {
	src := &person{Id: 1, Name: "Tim"}
	tls := NewThreadLocalWithInitial(nil)
	tls2 := NewThreadLocalWithInitial(func() Any {
		return nil
	})
	tls3 := NewThreadLocalWithInitial(func() Any {
		return src
	})
	tls4 := NewThreadLocalWithInitial(func() Any {
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

type person struct {
	Id   int
	Name string
}

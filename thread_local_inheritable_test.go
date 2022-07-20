package routine

import (
	"math"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInheritableThreadLocal_Index(t *testing.T) {
	tls := NewInheritableThreadLocal()
	assert.GreaterOrEqual(t, tls.(*inheritableThreadLocal).index, 0)
	tls2 := NewInheritableThreadLocalWithInitial(func() any {
		return "Hello"
	})
	assert.Greater(t, tls2.(*inheritableThreadLocal).index, tls.(*inheritableThreadLocal).index)
}

func TestInheritableThreadLocal_NextIndex(t *testing.T) {
	backup := inheritableThreadLocalIndex
	defer func() {
		inheritableThreadLocalIndex = backup
	}()
	//
	inheritableThreadLocalIndex = math.MaxInt32
	assert.Panics(t, func() {
		nextInheritableThreadLocalIndex()
	})
	assert.Equal(t, math.MaxInt32, int(inheritableThreadLocalIndex))
}

func TestInheritableThreadLocal_Common(t *testing.T) {
	tls := NewInheritableThreadLocal()
	tls2 := NewInheritableThreadLocal()
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
			assert.Equal(t, 2, tls.Get())
			assert.Equal(t, "!", tls2.Get())
			wg.Done()
		})
	}
	wg.Wait()
	assert.Equal(t, 2, tls.Get())
	assert.Equal(t, "!", tls2.Get())
}

func TestInheritableThreadLocal_Mixed(t *testing.T) {
	tls := NewInheritableThreadLocal()
	tls2 := NewInheritableThreadLocalWithInitial(func() any {
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
			assert.Equal(t, 2, tls.Get())
			assert.Equal(t, "!", tls2.Get())
			wg.Done()
		})
	}
	wg.Wait()
	assert.Equal(t, 2, tls.Get())
	assert.Equal(t, "!", tls2.Get())
}

func TestInheritableThreadLocal_WithInitial(t *testing.T) {
	src := &person{Id: 1, Name: "Tim"}
	tls := NewInheritableThreadLocalWithInitial(nil)
	tls2 := NewInheritableThreadLocalWithInitial(func() any {
		return nil
	})
	tls3 := NewInheritableThreadLocalWithInitial(func() any {
		return src
	})
	tls4 := NewInheritableThreadLocalWithInitial(func() any {
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

func TestInheritableThreadLocal_CrossCoroutine(t *testing.T) {
	tls := NewInheritableThreadLocal()
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
		assert.Equal(t, "Hello", tls.Get())
		finishWait.Done()
	})
	tls.Remove()      //remove in parent goroutine should not affect child goroutine
	subWait.Done()    //allow sub goroutine run
	subWait.Done()    //allow sub goroutine run
	finishWait.Wait() //wait sub goroutine done
	finishWait.Wait() //wait sub goroutine done
}

func TestInheritableThreadLocal_CreateBatch(t *testing.T) {
	const count = 128
	tlsList := make([]ThreadLocal, count)
	for i := 0; i < count; i++ {
		value := i
		tlsList[i] = NewInheritableThreadLocalWithInitial(func() any { return value })
	}
	for i := 0; i < count; i++ {
		assert.Equal(t, i, tlsList[i].Get())
	}
}

func TestInheritableThreadLocal_Copy(t *testing.T) {
	tls := NewInheritableThreadLocalWithInitial(func() any {
		return &person{Id: 1, Name: "Tim"}
	})
	tls2 := NewInheritableThreadLocalWithInitial(func() any {
		return person{Id: 2, Name: "Andy"}
	})

	p1 := tls.Get().(*person)
	assert.Equal(t, 1, p1.Id)
	assert.Equal(t, "Tim", p1.Name)
	p2 := tls2.Get().(person)
	assert.Equal(t, 2, p2.Id)
	assert.Equal(t, "Andy", p2.Name)
	//
	fut := GoWait(func(token CancelToken) {
		p3 := tls.Get().(*person)
		assert.Same(t, p1, p3)
		assert.Equal(t, 1, p3.Id)
		assert.Equal(t, "Tim", p1.Name)
		p4 := tls2.Get().(person)
		assert.NotSame(t, &p2, &p4)
		assert.Equal(t, p2, p4)
		assert.Equal(t, 2, p4.Id)
		assert.Equal(t, "Andy", p4.Name)
		//
		p3.Name = "Tim2"
		p4.Name = "Andy2"
	})
	fut.Get()
	//
	p5 := tls.Get().(*person)
	assert.Same(t, p1, p5)
	assert.Equal(t, 1, p5.Id)
	assert.Equal(t, "Tim2", p5.Name)
	p6 := tls2.Get().(person)
	assert.NotSame(t, &p2, &p6)
	assert.Equal(t, p2, p6)
	assert.Equal(t, 2, p6.Id)
	assert.Equal(t, "Andy", p6.Name)
}

func TestInheritableThreadLocal_Cloneable(t *testing.T) {
	tls := NewInheritableThreadLocalWithInitial(func() any {
		return &personCloneable{Id: 1, Name: "Tim"}
	})
	tls2 := NewInheritableThreadLocalWithInitial(func() any {
		return personCloneable{Id: 2, Name: "Andy"}
	})

	p1 := tls.Get().(*personCloneable)
	assert.Equal(t, 1, p1.Id)
	assert.Equal(t, "Tim", p1.Name)
	p2 := tls2.Get().(personCloneable)
	assert.Equal(t, 2, p2.Id)
	assert.Equal(t, "Andy", p2.Name)
	//
	fut := GoWait(func(token CancelToken) {
		p3 := tls.Get().(*personCloneable) //p3 is clone from p1
		assert.NotSame(t, p1, p3)
		assert.Equal(t, 1, p3.Id)
		assert.Equal(t, "Tim", p1.Name)
		p4 := tls2.Get().(personCloneable)
		assert.NotSame(t, &p2, &p4)
		assert.Equal(t, p2, p4)
		assert.Equal(t, 2, p4.Id)
		assert.Equal(t, "Andy", p4.Name)
		//
		p3.Name = "Tim2"
		p4.Name = "Andy2"
	})
	fut.Get()
	//
	p5 := tls.Get().(*personCloneable)
	assert.Same(t, p1, p5)
	assert.Equal(t, 1, p5.Id)
	assert.Equal(t, "Tim", p5.Name)
	p6 := tls2.Get().(personCloneable)
	assert.NotSame(t, &p2, &p6)
	assert.Equal(t, p2, p6)
	assert.Equal(t, 2, p6.Id)
	assert.Equal(t, "Andy", p6.Name)
}

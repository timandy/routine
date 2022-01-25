package routine

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestInheritableThreadLocal_Id(t *testing.T) {
	tls := NewInheritableThreadLocal()
	assert.GreaterOrEqual(t, tls.Id(), 0)
	tls2 := NewInheritableThreadLocalWithInitial(func() Any {
		return "Hello"
	})
	assert.Greater(t, tls2.Id(), tls.Id())
}

func TestInheritableThreadLocal(t *testing.T) {
	tls := NewInheritableThreadLocal()
	tls2 := NewInheritableThreadLocal()
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

func TestInheritableThreadLocalMixed(t *testing.T) {
	tls := NewInheritableThreadLocal()
	tls2 := NewInheritableThreadLocalWithInitial(func() Any {
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

func TestInheritableThreadLocalWithInitial(t *testing.T) {
	src := &person{Id: 1, Name: "Tim"}
	tls := NewInheritableThreadLocalWithInitial(nil)
	tls2 := NewInheritableThreadLocalWithInitial(func() Any {
		return nil
	})
	tls3 := NewInheritableThreadLocalWithInitial(func() Any {
		return src
	})
	tls4 := NewInheritableThreadLocalWithInitial(func() Any {
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

func TestInheritableThreadLocalCopy(t *testing.T) {
	tls := NewInheritableThreadLocalWithInitial(func() Any {
		return &person{Id: 1, Name: "Tim"}
	})
	tls2 := NewInheritableThreadLocalWithInitial(func() Any {
		return person{Id: 2, Name: "Andy"}
	})

	p1 := tls.Get().(*person)
	assert.Equal(t, 1, p1.Id)
	assert.Equal(t, "Tim", p1.Name)
	p2 := tls2.Get().(person)
	assert.Equal(t, 2, p2.Id)
	assert.Equal(t, "Andy", p2.Name)
	//
	fea := GoWait(func() {
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
	fea.Get()
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

func TestInheritableThreadLocalCloneable(t *testing.T) {
	tls := NewInheritableThreadLocalWithInitial(func() Any {
		return &personCloneable{Id: 1, Name: "Tim"}
	})
	tls2 := NewInheritableThreadLocalWithInitial(func() Any {
		return personCloneable{Id: 2, Name: "Andy"}
	})

	p1 := tls.Get().(*personCloneable)
	assert.Equal(t, 1, p1.Id)
	assert.Equal(t, "Tim", p1.Name)
	p2 := tls2.Get().(personCloneable)
	assert.Equal(t, 2, p2.Id)
	assert.Equal(t, "Andy", p2.Name)
	//
	fea := GoWait(func() {
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
	fea.Get()
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

package routine

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func init() {
	gCInterval = time.Millisecond * 50 // for faster test
}

func TestThreadGC(t *testing.T) {
	threadLocal := NewThreadLocal()
	threadLocal2 := NewThreadLocal()
	threadLocal3 := NewThreadLocalWithInitial(func() Any {
		return "Hello"
	})
	threadLocal4 := NewInheritableThreadLocal()
	threadLocal5 := NewInheritableThreadLocalWithInitial(func() Any {
		return 1234
	})

	// use ThreadLocal in multi goroutines
	for i := 0; i < 10; i++ {
		waiter := &sync.WaitGroup{}
		for i := 0; i < 1000; i++ {
			waiter.Add(1)
			Go(func() {
				threadLocal.Set("hello world")
				threadLocal2.Set(true)
				threadLocal3.Set(&threadLocal3)
				threadLocal4.Set(rand.Int())
				threadLocal5.Set(time.Now())
				waiter.Done()
			})
		}
		assert.True(t, gcRunning(), "#%v, timer may not running!", i)
		waiter.Wait()
		// wait for a while
		time.Sleep(gCInterval + time.Second)
		assert.False(t, gcRunning(), "#%v, timer not stopped?", i)
		gMap := globalMap.Load().(map[int64]*thread)
		assert.Equal(t, 0, len(gMap), "#%v, gMap not empty - %d", i, len(gMap))
	}
}

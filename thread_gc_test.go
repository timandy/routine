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
	tls := NewThreadLocal()
	tls2 := NewThreadLocal()
	tls3 := NewThreadLocalWithInitial(func() Any {
		return "Hello"
	})
	tls4 := NewInheritableThreadLocal()
	tls5 := NewInheritableThreadLocalWithInitial(func() Any {
		return 1234
	})

	// use ThreadLocal in multi goroutines
	for i := 0; i < 10; i++ {
		waiter := &sync.WaitGroup{}
		for j := 0; j < 1000; j++ {
			waiter.Add(1)
			Go(func() {
				tls.Set("hello world")
				tls2.Set(true)
				tls3.Set(&tls3)
				tls4.Set(rand.Int())
				tls5.Set(time.Now())
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

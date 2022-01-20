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
	const loopTimes = 10
	const concurrency = 1000
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
	gcCnt := 0
	for i := 0; i < loopTimes; i++ {
		wg := &sync.WaitGroup{}
		wg.Add(concurrency)
		for j := 0; j < concurrency; j++ {
			Go(func() {
				tls.Set("hello world")
				tls2.Set(true)
				tls3.Set(&tls3)
				tls4.Set(rand.Int())
				tls5.Set(time.Now())
				wg.Done()
			})
		}
		if gcRunning() {
			gcCnt++
		}
		wg.Wait()
		// wait for a while
		time.Sleep(gCInterval + 200*time.Millisecond)
		assert.False(t, gcRunning(), "#%v, gcTimer not stopped!", i)
		gMap := globalMap.Load().(map[int64]*thread)
		assert.Equal(t, 0, len(gMap), "#%v, gMap(len=%d) not empty after gc!", i, len(gMap))
	}
	if gcCnt != loopTimes {
		t.Logf("[WARNING] gcTimer running count is: %v!", gcCnt)
	}
	assert.Greater(t, gcCnt, loopTimes/2, "gcTimer not running!")
}

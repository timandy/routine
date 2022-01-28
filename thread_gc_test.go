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
		if hasAnyGcTimer() {
			gcCnt++
		}
		wg.Wait()
		// wait for a while
		time.Sleep(gCInterval + 200*time.Millisecond)
		assert.False(t, hasAnyGcTimer(), "#%v, gcTimer not stopped!", i)
		threadsCnt := countAllThreads()
		assert.Equal(t, 0, threadsCnt, "#%v, globalMap(len=%d) not empty after gc!", i, threadsCnt)
	}
	if gcCnt != loopTimes {
		t.Logf("[WARNING] gcTimer running count is: %v!", gcCnt)
	}
	assert.Greater(t, gcCnt, loopTimes/2, "gcTimer not running!")
}

func hasAnyGcTimer() bool {
	for idx := 0; idx < segmentSize; idx++ {
		if hasGcTimer(idx) {
			return true
		}
	}
	return false
}

func hasGcTimer(idx int) bool {
	segmentLock := globalMapLock[idx]
	segmentLock.Lock()
	defer segmentLock.Unlock()
	return gcTimer[idx] != nil
}

func countAllThreads() int {
	count := 0
	for idx := 0; idx < segmentSize; idx++ {
		count += countThreads(idx)
	}
	return count
}

func countThreads(idx int) int {
	segmentLock := globalMapLock[idx]
	segmentLock.Lock()
	defer segmentLock.Unlock()
	return len(globalMap[idx].Load().(map[int64]*thread))
}

package routine

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"runtime"
	"runtime/pprof"
	"sync"
	"testing"
	"time"
)

func TestCurrentThread(t *testing.T) {
	assert.NotNil(t, currentThread(true))
	assert.Same(t, currentThread(true), currentThread(true))
}

func TestPProf(t *testing.T) {
	const concurrency = 10
	const loopTimes = 10
	tls := NewThreadLocal()
	tls.Set("你好")
	wg := &sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		tmp := i
		go func() {
			for j := 0; j < loopTimes; j++ {
				time.Sleep(100 * time.Millisecond)
				tls.Set(tmp)
				assert.Equal(t, tmp, tls.Get())
				pprof.Do(context.Background(), pprof.Labels("key", "value"), func(ctx context.Context) {
					assert.Nil(t, currentThread(false))
					assert.Nil(t, tls.Get())
					tls.Set("hi")
					//
					label, find := pprof.Label(ctx, "key")
					assert.True(t, find)
					assert.Equal(t, "value", label)
					//
					assert.Equal(t, "hi", tls.Get())
					//
					label2, find2 := pprof.Label(ctx, "key")
					assert.True(t, find2)
					assert.Equal(t, "value", label2)
				})
				assert.Nil(t, tls.Get())
			}
			wg.Done()
		}()
	}
	assert.Nil(t, pprof.StartCPUProfile(&bytes.Buffer{}))
	wg.Wait()
	pprof.StopCPUProfile()
	assert.Equal(t, "你好", tls.Get())
}

func TestThreadGC(t *testing.T) {
	const allocSize = 10_000_000
	tls := NewThreadLocal()
	tls2 := NewInheritableThreadLocal()
	allocWait := &sync.WaitGroup{}
	allocWait.Add(1)
	gatherWait := &sync.WaitGroup{}
	gatherWait.Add(1)
	gcWait := &sync.WaitGroup{}
	gcWait.Add(1)
	//=========Init
	heapInit, numInit := getMemStats()
	printMemStats("Init", heapInit, numInit)
	//
	fea := GoWait(func() {
		tls.Set(make([]byte, allocSize))
		tls2.Set(make([]byte, allocSize))
		go func() {
			gcWait.Wait()
		}()
		fea2 := GoWaitResult(func() Any {
			return 1
		})
		assert.Equal(t, 1, fea2.Get())
		allocWait.Done()  //alloc ok, release main thread
		gatherWait.Wait() //wait gather heap info
	})
	//=========Alloc
	allocWait.Wait() //wait alloc done
	heapAlloc, numAlloc := getMemStats()
	printMemStats("Alloc", heapAlloc, numAlloc)
	assert.Greater(t, heapAlloc, heapInit+allocSize*2*0.9)
	assert.Greater(t, numAlloc, numInit)
	//=========GC
	gatherWait.Done() //gather ok, release sub thread
	fea.Get()         //wait sub thread finish
	time.Sleep(500 * time.Millisecond)
	heapGC, numGC := getMemStats()
	printMemStats("AfterGC", heapGC, numGC)
	gcWait.Done()
	//=========Summary
	heapRelease := heapAlloc - heapGC
	numRelease := numAlloc - numGC
	printMemStats("Summary", heapRelease, numRelease)
	assert.Greater(t, int(heapRelease), int(allocSize*2*0.9))
	assert.Equal(t, 1, numRelease)
}

func getMemStats() (uint64, int) {
	stats := runtime.MemStats{}
	runtime.GC()
	runtime.ReadMemStats(&stats)
	return stats.HeapAlloc, runtime.NumGoroutine()
}

func printMemStats(section string, heapAlloc uint64, numGoroutine int) {
	//fmt.Printf("%v\n", section)
	//fmt.Printf("HeapAlloc    = %v\n", heapAlloc)
	//fmt.Printf("NumGoroutine = %v\n", numGoroutine)
	//fmt.Printf("===\n")
}

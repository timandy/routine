package routine

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestCurrentThread(t *testing.T) {
	assert.NotNil(t, currentThread())
	assert.Same(t, currentThread(), currentThread())
}

func TestThreadGC(t *testing.T) {
	const allocSize = 10_000_000
	tls := NewThreadLocal()
	allocWait := &sync.WaitGroup{}
	allocWait.Add(1)
	gatherWait := &sync.WaitGroup{}
	gatherWait.Add(1)
	//=========Init
	heapInit, numInit := getMemStats()
	printMemStats("Init", heapInit, numInit)
	//
	fea := GoWait(func() {
		tls.Set(make([]byte, allocSize))
		allocWait.Done()  //alloc ok, release main thread
		gatherWait.Wait() //wait gather heap info
	})
	//=========Alloc
	allocWait.Wait() //wait alloc done
	heapAlloc, numAlloc := getMemStats()
	printMemStats("Alloc", heapAlloc, numAlloc)
	assert.Greater(t, heapAlloc, heapInit)
	assert.Greater(t, numAlloc, numInit)
	//=========GC
	gatherWait.Done() //gather ok, release sub thread
	fea.Get()         //wait sub thread finish
	time.Sleep(time.Millisecond * 500)
	heapGC, numGC := getMemStats()
	printMemStats("AfterGC", heapGC, numGC)
	//=========Summary
	heapRelease := heapAlloc - heapGC
	numRelease := numAlloc - numGC
	printMemStats("Summary", heapRelease, numRelease)
	assert.Greater(t, int(heapRelease), int(allocSize*0.9))
	assert.Equal(t, 1, numRelease)
}

func getMemStats() (uint64, int) {
	var stats runtime.MemStats
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

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

func TestThreadLocal(t *testing.T) {
	s := NewThreadLocal()
	s2 := NewThreadLocal()

	for i := 0; i < 100; i++ {
		src := "hello"
		s.Set(src)
		p := s.Get()
		assert.Equal(t, src, p.(string))
		//
		src2 := "world"
		s2.Set(src2)
		p2 := s2.Get()
		assert.Equal(t, src2, p2.(string))
	}

	for i := 0; i < 1000; i++ {
		num := rand.Int()
		s.Set(num)
		num2 := s.Get()
		assert.Equal(t, num, num2.(int))
	}

	v := s.Remove()
	assert.NotNil(t, v)

	Clear()
	vv1 := s.Get()
	assert.Nil(t, vv1)
	//
	vv2 := s2.Get()
	assert.Nil(t, vv2)
}

func TestThreadLocalConcurrency(t *testing.T) {
	const concurrency = 1000
	const loopTimes = 1000

	s := NewThreadLocal()
	s2 := NewThreadLocal()

	waiter := &sync.WaitGroup{}
	waiter.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			v := rand.Uint64()
			v2 := rand.Uint64()
			assert.True(t, v != 0)
			assert.True(t, v2 != 0)
			for i := 0; i < loopTimes; i++ {
				s.Set(v)
				tmp := s.Get()
				assert.True(t, tmp.(uint64) == v)
				//
				s2.Set(v2)
				tmp2 := s2.Get()
				assert.True(t, tmp2.(uint64) == v2)
			}
			waiter.Done()
		}()
	}
	waiter.Wait()
}

func TestThreadLocalGC(t *testing.T) {
	s1 := NewThreadLocal()
	s2 := NewThreadLocal()
	s3 := NewThreadLocal()
	s4 := NewThreadLocal()
	s5 := NewThreadLocal()

	// use ThreadLocal in multi goroutines
	gcRunningCnt := 0
	for i := 0; i < 10; i++ {
		waiter := &sync.WaitGroup{}
		for i := 0; i < 1000; i++ {
			waiter.Add(1)
			go func() {
				s1.Set("hello world")
				s2.Set(true)
				s3.Set(&s3)
				s4.Set(rand.Int())
				s5.Set(time.Now())
				waiter.Done()
			}()
		}
		if gcRunning() {
			// record gc running count
			gcRunningCnt++
		}
		waiter.Wait()
		// wait for a while
		time.Sleep(gCInterval + time.Second)
		assert.False(t, gcRunning(), "#%v, timer not stoped?", i)
		storeMap := globalMap.Load().(map[int64]*threadLocalMap)
		assert.Equal(t, 0, len(storeMap), "#%v, storeMap not empty - %d", i, len(storeMap))
	}
	assert.True(t, gcRunningCnt > 0, "gc timer may not running!")
}

// BenchmarkThreadLocal-8   	10183446	        99.28 ns/op	       8 B/op	       0 allocs/op
func BenchmarkThreadLocal(b *testing.B) {
	threadLocalCount := 100
	threadLocals := make([]ThreadLocal, threadLocalCount)
	for i := 0; i < threadLocalCount; i++ {
		threadLocals[i] = NewThreadLocal()
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % threadLocalCount
		threadLocal := threadLocals[index]
		threadLocal.Set(i)
		if threadLocal.Get() != i {
			b.Fail()
		}
		threadLocal.Remove()
	}
}

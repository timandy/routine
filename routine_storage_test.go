package routine

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func init() {
	storageGCInterval = time.Millisecond * 50 // for faster test
}

func TestStorage(t *testing.T) {
	s := NewLocalStorage()
	s2 := NewLocalStorage()

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

func TestStorageConcurrency(t *testing.T) {
	const concurrency = 1000
	const loopTimes = 1000

	s := NewLocalStorage()
	s2 := NewLocalStorage()

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

func TestStorageGC(t *testing.T) {
	s1 := NewLocalStorage()
	s2 := NewLocalStorage()
	s3 := NewLocalStorage()
	s4 := NewLocalStorage()
	s5 := NewLocalStorage()

	// use LocalStorage in multi goroutines
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
		assert.True(t, gcRunning(), "#%v, timer not running?", i)
		waiter.Wait()
		// wait for a while
		time.Sleep(storageGCInterval + time.Second)
		assert.False(t, gcRunning(), "#%v, timer not stoped?", i)
		storeMap := storages.Load().(map[int64]*store)
		assert.Equal(t, 0, len(storeMap), "#%v, storeMap not empty - %d", i, len(storeMap))
	}
}

// BenchmarkLoadCurrentStore-12    	 9630090	       118.2 ns/op	      16 B/op	       1 allocs/op
func BenchmarkStorage(b *testing.B) {
	s := NewLocalStorage()
	variable := "hello world"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Get()
		s.Set(variable)
		s.Remove()
	}
}

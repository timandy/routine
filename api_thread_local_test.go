package routine

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
)

const (
	concurrency = 1000
	loopTimes   = 200
)

func TestSupplier(t *testing.T) {
	var fun Supplier
	fun = func() interface{} {
		return "Hello"
	}
	assert.Equal(t, "Hello", fun())
}

//===

func TestThreadLocal_New(t *testing.T) {
	threadLocal := NewThreadLocal()
	threadLocal.Set("Hello")
	assert.Equal(t, "Hello", threadLocal.Get())
	//
	threadLocal2 := NewThreadLocal()
	assert.Equal(t, "Hello", threadLocal.Get())
	threadLocal2.Set(22)
	assert.Equal(t, 22, threadLocal2.Get())
	//
	threadLocal2.Set(33)
	assert.Equal(t, 33, threadLocal2.Get())
	//
	feature := GoWait(func() {
		assert.Nil(t, threadLocal.Get())
		assert.Nil(t, threadLocal2.Get())
	})
	feature.Get()
}

func TestThreadLocal_Multi(t *testing.T) {
	threadLocal := NewThreadLocal()
	threadLocal2 := NewThreadLocal()
	threadLocal.Set("Hello")
	threadLocal2.Set(22)
	assert.Equal(t, 22, threadLocal2.Get())
	assert.Equal(t, "Hello", threadLocal.Get())
	//
	threadLocal2.Set(33)
	assert.Equal(t, 33, threadLocal2.Get())
	//
	feature := GoWait(func() {
		assert.Nil(t, threadLocal.Get())
		assert.Nil(t, threadLocal2.Get())
	})
	feature.Get()
}

func TestThreadLocal_Concurrency(t *testing.T) {
	threadLocal := NewThreadLocal()
	threadLocal2 := NewThreadLocal()
	//
	threadLocal2.Set(33)
	assert.Equal(t, 33, threadLocal2.Get())
	//
	waiter := &sync.WaitGroup{}
	waiter.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		assert.Nil(t, threadLocal.Get())
		assert.Equal(t, 33, threadLocal2.Get())
		Go(func() {
			assert.Nil(t, threadLocal.Get())
			assert.Nil(t, threadLocal2.Get())
			v := rand.Uint64()
			v2 := rand.Uint64()
			for i := 0; i < loopTimes; i++ {
				threadLocal.Set(v)
				tmp := threadLocal.Get()
				assert.Equal(t, v, tmp.(uint64))
				//
				threadLocal2.Set(v2)
				tmp2 := threadLocal2.Get()
				assert.Equal(t, v2, tmp2.(uint64))
			}
			waiter.Done()
		})
	}
	waiter.Wait()
	//
	feature := GoWait(func() {
		assert.Nil(t, threadLocal.Get())
		assert.Nil(t, threadLocal2.Get())
	})
	feature.Get()
}

//===

func TestThreadLocalWithInitial_New(t *testing.T) {
	threadLocal := NewThreadLocalWithInitial(func() Any {
		return "Hello"
	})
	assert.Equal(t, "Hello", threadLocal.Get())
	//
	threadLocal2 := NewThreadLocalWithInitial(func() Any {
		return 22
	})
	assert.Equal(t, "Hello", threadLocal.Get())
	assert.Equal(t, 22, threadLocal2.Get())
	//
	threadLocal2.Set(33)
	assert.Equal(t, 33, threadLocal2.Get())
	//
	feature := GoWait(func() {
		assert.Equal(t, "Hello", threadLocal.Get())
		assert.Equal(t, 22, threadLocal2.Get())
	})
	feature.Get()
}

func TestThreadLocalWithInitial_Multi(t *testing.T) {
	threadLocal := NewThreadLocalWithInitial(func() Any {
		return "Hello"
	})
	threadLocal2 := NewThreadLocalWithInitial(func() Any {
		return 22
	})
	threadLocal.Set("Hello")
	threadLocal2.Set(22)
	assert.Equal(t, 22, threadLocal2.Get())
	assert.Equal(t, "Hello", threadLocal.Get())
	//
	threadLocal2.Set(33)
	assert.Equal(t, 33, threadLocal2.Get())
	//
	feature := GoWait(func() {
		assert.Equal(t, "Hello", threadLocal.Get())
		assert.Equal(t, 22, threadLocal2.Get())
	})
	feature.Get()
}

func TestThreadLocalWithInitial_Concurrency(t *testing.T) {
	threadLocal := NewThreadLocalWithInitial(func() Any {
		return "Hello"
	})
	threadLocal2 := NewThreadLocalWithInitial(func() Any {
		return 22
	})
	//
	threadLocal2.Set(33)
	assert.Equal(t, 33, threadLocal2.Get())
	//
	waiter := &sync.WaitGroup{}
	waiter.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		assert.Equal(t, "Hello", threadLocal.Get())
		assert.Equal(t, 33, threadLocal2.Get())
		Go(func() {
			assert.Equal(t, "Hello", threadLocal.Get())
			assert.Equal(t, 22, threadLocal2.Get())
			v := rand.Uint64()
			v2 := rand.Uint64()
			for i := 0; i < loopTimes; i++ {
				threadLocal.Set(v)
				tmp := threadLocal.Get()
				assert.Equal(t, v, tmp.(uint64))
				//
				threadLocal2.Set(v2)
				tmp2 := threadLocal2.Get()
				assert.Equal(t, v2, tmp2.(uint64))
			}
			waiter.Done()
		})
	}
	waiter.Wait()
	//
	feature := GoWait(func() {
		assert.Equal(t, "Hello", threadLocal.Get())
		assert.Equal(t, 22, threadLocal2.Get())
	})
	feature.Get()
}

//===

func TestInheritableThreadLocal_New(t *testing.T) {
	threadLocal := NewInheritableThreadLocal()
	threadLocal.Set("Hello")
	assert.Equal(t, "Hello", threadLocal.Get())
	//
	threadLocal2 := NewInheritableThreadLocal()
	assert.Equal(t, "Hello", threadLocal.Get())
	threadLocal2.Set(22)
	assert.Equal(t, 22, threadLocal2.Get())
	//
	threadLocal2.Set(33)
	assert.Equal(t, 33, threadLocal2.Get())
	//
	feature := GoWait(func() {
		assert.Equal(t, "Hello", threadLocal.Get())
		assert.Equal(t, 33, threadLocal2.Get())
	})
	feature.Get()
}

func TestInheritableThreadLocal_Multi(t *testing.T) {
	threadLocal := NewInheritableThreadLocal()
	threadLocal2 := NewInheritableThreadLocal()
	threadLocal.Set("Hello")
	threadLocal2.Set(22)
	assert.Equal(t, 22, threadLocal2.Get())
	assert.Equal(t, "Hello", threadLocal.Get())
	//
	threadLocal2.Set(33)
	assert.Equal(t, 33, threadLocal2.Get())
	//
	feature := GoWait(func() {
		assert.Equal(t, "Hello", threadLocal.Get())
		assert.Equal(t, 33, threadLocal2.Get())
	})
	feature.Get()
}

func TestInheritableThreadLocal_Concurrency(t *testing.T) {
	threadLocal := NewInheritableThreadLocal()
	threadLocal2 := NewInheritableThreadLocal()
	//
	threadLocal2.Set(33)
	assert.Equal(t, 33, threadLocal2.Get())
	//
	waiter := &sync.WaitGroup{}
	waiter.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		assert.Nil(t, threadLocal.Get())
		assert.Equal(t, 33, threadLocal2.Get())
		Go(func() {
			assert.Nil(t, threadLocal.Get())
			assert.Equal(t, 33, threadLocal2.Get())
			v := rand.Uint64()
			v2 := rand.Uint64()
			for i := 0; i < loopTimes; i++ {
				threadLocal.Set(v)
				tmp := threadLocal.Get()
				assert.Equal(t, v, tmp.(uint64))
				//
				threadLocal2.Set(v2)
				tmp2 := threadLocal2.Get()
				assert.Equal(t, v2, tmp2.(uint64))
			}
			waiter.Done()
		})
	}
	waiter.Wait()
	//
	feature := GoWait(func() {
		assert.Nil(t, threadLocal.Get())
		assert.Equal(t, 33, threadLocal2.Get())
	})
	feature.Get()
}

//===

func TestInheritableThreadLocalWithInitial_New(t *testing.T) {
	threadLocal := NewInheritableThreadLocalWithInitial(func() Any {
		return "Hello"
	})
	assert.Equal(t, "Hello", threadLocal.Get())
	//
	threadLocal2 := NewInheritableThreadLocalWithInitial(func() Any {
		return 22
	})
	assert.Equal(t, "Hello", threadLocal.Get())
	assert.Equal(t, 22, threadLocal2.Get())
	//
	threadLocal2.Set(33)
	assert.Equal(t, 33, threadLocal2.Get())
	//
	feature := GoWait(func() {
		assert.Equal(t, "Hello", threadLocal.Get())
		assert.Equal(t, 33, threadLocal2.Get())
	})
	feature.Get()
}

func TestInheritableThreadLocalWithInitial_Multi(t *testing.T) {
	threadLocal := NewInheritableThreadLocalWithInitial(func() Any {
		return "Hello"
	})
	threadLocal2 := NewInheritableThreadLocalWithInitial(func() Any {
		return 22
	})
	threadLocal.Set("Hello")
	threadLocal2.Set(22)
	assert.Equal(t, 22, threadLocal2.Get())
	assert.Equal(t, "Hello", threadLocal.Get())
	//
	threadLocal2.Set(33)
	assert.Equal(t, 33, threadLocal2.Get())
	//
	feature := GoWait(func() {
		assert.Equal(t, "Hello", threadLocal.Get())
		assert.Equal(t, 33, threadLocal2.Get())
	})
	feature.Get()
}

func TestInheritableThreadLocalWithInitial_Concurrency(t *testing.T) {
	threadLocal := NewInheritableThreadLocalWithInitial(func() Any {
		return "Hello"
	})
	threadLocal2 := NewInheritableThreadLocalWithInitial(func() Any {
		return 22
	})
	//
	threadLocal2.Set(33)
	assert.Equal(t, 33, threadLocal2.Get())
	//
	waiter := &sync.WaitGroup{}
	waiter.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		assert.Equal(t, "Hello", threadLocal.Get())
		assert.Equal(t, 33, threadLocal2.Get())
		Go(func() {
			assert.Equal(t, "Hello", threadLocal.Get())
			assert.Equal(t, 33, threadLocal2.Get())
			v := rand.Uint64()
			v2 := rand.Uint64()
			for i := 0; i < loopTimes; i++ {
				threadLocal.Set(v)
				tmp := threadLocal.Get()
				assert.Equal(t, v, tmp.(uint64))
				//
				threadLocal2.Set(v2)
				tmp2 := threadLocal2.Get()
				assert.Equal(t, v2, tmp2.(uint64))
			}
			waiter.Done()
		})
	}
	waiter.Wait()
	//
	feature := GoWait(func() {
		assert.Equal(t, "Hello", threadLocal.Get())
		assert.Equal(t, 33, threadLocal2.Get())
	})
	feature.Get()
}

//===

// BenchmarkThreadLocal-4                           7792273               159.2 ns/op             8 B/op          0 allocs/op
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
		initValue := threadLocal.Get()
		if initValue != nil {
			b.Fail()
		}
		threadLocal.Set(i)
		if threadLocal.Get() != i {
			b.Fail()
		}
		threadLocal.Remove()
	}
}

// BenchmarkThreadLocalWithInitial-4                7868341               151.1 ns/op             8 B/op          0 allocs/op
func BenchmarkThreadLocalWithInitial(b *testing.B) {
	threadLocalCount := 100
	threadLocals := make([]ThreadLocal, threadLocalCount)
	for i := 0; i < threadLocalCount; i++ {
		index := i
		threadLocals[i] = NewThreadLocalWithInitial(func() Any {
			return index
		})
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % threadLocalCount
		threadLocal := threadLocals[index]
		initValue := threadLocal.Get()
		if initValue != index {
			b.Fail()
		}
		threadLocal.Set(i)
		if threadLocal.Get() != i {
			b.Fail()
		}
		threadLocal.Remove()
	}
}

// BenchmarkInheritableThreadLocal-4                8228482               150.8 ns/op             8 B/op          0 allocs/op
func BenchmarkInheritableThreadLocal(b *testing.B) {
	threadLocalCount := 100
	threadLocals := make([]ThreadLocal, threadLocalCount)
	for i := 0; i < threadLocalCount; i++ {
		threadLocals[i] = NewInheritableThreadLocal()
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % threadLocalCount
		threadLocal := threadLocals[index]
		initValue := threadLocal.Get()
		if initValue != nil {
			b.Fail()
		}
		threadLocal.Set(i)
		if threadLocal.Get() != i {
			b.Fail()
		}
		threadLocal.Remove()
	}
}

// BenchmarkInheritableThreadLocalWithInitial-4     7407096               159.1 ns/op             8 B/op          0 allocs/op
func BenchmarkInheritableThreadLocalWithInitial(b *testing.B) {
	threadLocalCount := 100
	threadLocals := make([]ThreadLocal, threadLocalCount)
	for i := 0; i < threadLocalCount; i++ {
		index := i
		threadLocals[i] = NewInheritableThreadLocalWithInitial(func() Any {
			return index
		})
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % threadLocalCount
		threadLocal := threadLocals[index]
		initValue := threadLocal.Get()
		if initValue != index {
			b.Fail()
		}
		threadLocal.Set(i)
		if threadLocal.Get() != i {
			b.Fail()
		}
		threadLocal.Remove()
	}
}

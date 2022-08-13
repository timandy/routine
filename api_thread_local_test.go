package routine

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	concurrency = 1000
	loopTimes   = 200
)

func TestSupplier(t *testing.T) {
	var supplier Supplier = func() interface{} {
		return "Hello"
	}
	assert.Equal(t, "Hello", supplier())
	//
	var fun func() any = supplier
	assert.Equal(t, "Hello", fun())
}

//===

func TestNewThreadLocal_Single(t *testing.T) {
	tls := NewThreadLocal()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2 := NewThreadLocal()
	assert.Equal(t, "Hello", tls.Get())
	tls2.Set(22)
	assert.Equal(t, 22, tls2.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	fut := GoWait(func(token CancelToken) {
		assert.Nil(t, tls.Get())
		assert.Nil(t, tls2.Get())
	})
	fut.Get()
}

func TestNewThreadLocal_Multi(t *testing.T) {
	tls := NewThreadLocal()
	tls2 := NewThreadLocal()
	tls.Set("Hello")
	tls2.Set(22)
	assert.Equal(t, 22, tls2.Get())
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	fut := GoWait(func(token CancelToken) {
		assert.Nil(t, tls.Get())
		assert.Nil(t, tls2.Get())
	})
	fut.Get()
}

func TestNewThreadLocal_Concurrency(t *testing.T) {
	tls := NewThreadLocal()
	tls2 := NewThreadLocal()
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	wg := &sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		assert.Nil(t, tls.Get())
		assert.Equal(t, 33, tls2.Get())
		Go(func() {
			assert.Nil(t, tls.Get())
			assert.Nil(t, tls2.Get())
			v := rand.Uint64()
			v2 := rand.Uint64()
			for j := 0; j < loopTimes; j++ {
				tls.Set(v)
				tmp := tls.Get()
				assert.Equal(t, v, tmp.(uint64))
				//
				tls2.Set(v2)
				tmp2 := tls2.Get()
				assert.Equal(t, v2, tmp2.(uint64))
			}
			wg.Done()
		})
	}
	wg.Wait()
	//
	fut := GoWait(func(token CancelToken) {
		assert.Nil(t, tls.Get())
		assert.Nil(t, tls2.Get())
	})
	fut.Get()
}

//===

func TestNewThreadLocalWithInitial_Single(t *testing.T) {
	tls := NewThreadLocalWithInitial(func() any {
		return "Hello"
	})
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2 := NewThreadLocalWithInitial(func() any {
		return 22
	})
	assert.Equal(t, "Hello", tls.Get())
	assert.Equal(t, 22, tls2.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	fut := GoWait(func(token CancelToken) {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 22, tls2.Get())
	})
	fut.Get()
}

func TestNewThreadLocalWithInitial_Multi(t *testing.T) {
	tls := NewThreadLocalWithInitial(func() any {
		return "Hello"
	})
	tls2 := NewThreadLocalWithInitial(func() any {
		return 22
	})
	tls.Set("Hello")
	tls2.Set(22)
	assert.Equal(t, 22, tls2.Get())
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	fut := GoWait(func(token CancelToken) {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 22, tls2.Get())
	})
	fut.Get()
}

func TestNewThreadLocalWithInitial_Concurrency(t *testing.T) {
	tls := NewThreadLocalWithInitial(func() any {
		return "Hello"
	})
	tls2 := NewThreadLocalWithInitial(func() any {
		return 22
	})
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	wg := &sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 33, tls2.Get())
		Go(func() {
			assert.Equal(t, "Hello", tls.Get())
			assert.Equal(t, 22, tls2.Get())
			v := rand.Uint64()
			v2 := rand.Uint64()
			for j := 0; j < loopTimes; j++ {
				tls.Set(v)
				tmp := tls.Get()
				assert.Equal(t, v, tmp.(uint64))
				//
				tls2.Set(v2)
				tmp2 := tls2.Get()
				assert.Equal(t, v2, tmp2.(uint64))
			}
			wg.Done()
		})
	}
	wg.Wait()
	//
	fut := GoWait(func(token CancelToken) {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 22, tls2.Get())
	})
	fut.Get()
}

//===

func TestNewInheritableThreadLocal_Single(t *testing.T) {
	tls := NewInheritableThreadLocal()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2 := NewInheritableThreadLocal()
	assert.Equal(t, "Hello", tls.Get())
	tls2.Set(22)
	assert.Equal(t, 22, tls2.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	fut := GoWait(func(token CancelToken) {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 33, tls2.Get())
	})
	fut.Get()
}

func TestNewInheritableThreadLocal_Multi(t *testing.T) {
	tls := NewInheritableThreadLocal()
	tls2 := NewInheritableThreadLocal()
	tls.Set("Hello")
	tls2.Set(22)
	assert.Equal(t, 22, tls2.Get())
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	fut := GoWait(func(token CancelToken) {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 33, tls2.Get())
	})
	fut.Get()
}

func TestNewInheritableThreadLocal_Concurrency(t *testing.T) {
	tls := NewInheritableThreadLocal()
	tls2 := NewInheritableThreadLocal()
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	wg := &sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		assert.Nil(t, tls.Get())
		assert.Equal(t, 33, tls2.Get())
		Go(func() {
			assert.Nil(t, tls.Get())
			assert.Equal(t, 33, tls2.Get())
			v := rand.Uint64()
			v2 := rand.Uint64()
			for j := 0; j < loopTimes; j++ {
				tls.Set(v)
				tmp := tls.Get()
				assert.Equal(t, v, tmp.(uint64))
				//
				tls2.Set(v2)
				tmp2 := tls2.Get()
				assert.Equal(t, v2, tmp2.(uint64))
			}
			wg.Done()
		})
	}
	wg.Wait()
	//
	fut := GoWait(func(token CancelToken) {
		assert.Nil(t, tls.Get())
		assert.Equal(t, 33, tls2.Get())
	})
	fut.Get()
}

//===

func TestNewInheritableThreadLocalWithInitial_Single(t *testing.T) {
	tls := NewInheritableThreadLocalWithInitial(func() any {
		return "Hello"
	})
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2 := NewInheritableThreadLocalWithInitial(func() any {
		return 22
	})
	assert.Equal(t, "Hello", tls.Get())
	assert.Equal(t, 22, tls2.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	fut := GoWait(func(token CancelToken) {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 33, tls2.Get())
	})
	fut.Get()
}

func TestNewInheritableThreadLocalWithInitial_Multi(t *testing.T) {
	tls := NewInheritableThreadLocalWithInitial(func() any {
		return "Hello"
	})
	tls2 := NewInheritableThreadLocalWithInitial(func() any {
		return 22
	})
	tls.Set("Hello")
	tls2.Set(22)
	assert.Equal(t, 22, tls2.Get())
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	fut := GoWait(func(token CancelToken) {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 33, tls2.Get())
	})
	fut.Get()
}

func TestNewInheritableThreadLocalWithInitial_Concurrency(t *testing.T) {
	tls := NewInheritableThreadLocalWithInitial(func() any {
		return "Hello"
	})
	tls2 := NewInheritableThreadLocalWithInitial(func() any {
		return 22
	})
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	wg := &sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 33, tls2.Get())
		Go(func() {
			assert.Equal(t, "Hello", tls.Get())
			assert.Equal(t, 33, tls2.Get())
			v := rand.Uint64()
			v2 := rand.Uint64()
			for j := 0; j < loopTimes; j++ {
				tls.Set(v)
				tmp := tls.Get()
				assert.Equal(t, v, tmp.(uint64))
				//
				tls2.Set(v2)
				tmp2 := tls2.Get()
				assert.Equal(t, v2, tmp2.(uint64))
			}
			wg.Done()
		})
	}
	wg.Wait()
	//
	fut := GoWait(func(token CancelToken) {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 33, tls2.Get())
	})
	fut.Get()
}

//===

// BenchmarkThreadLocal-4                          16088140                74.48 ns/op            7 B/op          0 allocs/op
func BenchmarkThreadLocal(b *testing.B) {
	tlsCount := 100
	tlsSlice := make([]ThreadLocal, tlsCount)
	for i := 0; i < tlsCount; i++ {
		tlsSlice[i] = NewThreadLocal()
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % tlsCount
		tls := tlsSlice[index]
		initValue := tls.Get()
		if initValue != nil {
			b.Fail()
		}
		tls.Set(i)
		if tls.Get() != i {
			b.Fail()
		}
		tls.Remove()
	}
}

// BenchmarkThreadLocalWithInitial-4               15618451                77.03 ns/op            7 B/op          0 allocs/op
func BenchmarkThreadLocalWithInitial(b *testing.B) {
	tlsCount := 100
	tlsSlice := make([]ThreadLocal, tlsCount)
	for i := 0; i < tlsCount; i++ {
		index := i
		tlsSlice[i] = NewThreadLocalWithInitial(func() any {
			return index
		})
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % tlsCount
		tls := tlsSlice[index]
		initValue := tls.Get()
		if initValue != index {
			b.Fail()
		}
		tls.Set(i)
		if tls.Get() != i {
			b.Fail()
		}
		tls.Remove()
	}
}

// BenchmarkInheritableThreadLocal-4               16109587                73.17 ns/op            7 B/op          0 allocs/op
func BenchmarkInheritableThreadLocal(b *testing.B) {
	tlsCount := 100
	tlsSlice := make([]ThreadLocal, tlsCount)
	for i := 0; i < tlsCount; i++ {
		tlsSlice[i] = NewInheritableThreadLocal()
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % tlsCount
		tls := tlsSlice[index]
		initValue := tls.Get()
		if initValue != nil {
			b.Fail()
		}
		tls.Set(i)
		if tls.Get() != i {
			b.Fail()
		}
		tls.Remove()
	}
}

// BenchmarkInheritableThreadLocalWithInitial-4    14862778                78.77 ns/op            7 B/op          0 allocs/op
func BenchmarkInheritableThreadLocalWithInitial(b *testing.B) {
	tlsCount := 100
	tlsSlice := make([]ThreadLocal, tlsCount)
	for i := 0; i < tlsCount; i++ {
		index := i
		tlsSlice[i] = NewInheritableThreadLocalWithInitial(func() any {
			return index
		})
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % tlsCount
		tls := tlsSlice[index]
		initValue := tls.Get()
		if initValue != index {
			b.Fail()
		}
		tls.Set(i)
		if tls.Get() != i {
			b.Fail()
		}
		tls.Remove()
	}
}

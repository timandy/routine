package routine

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	concurrency = 500
	loopTimes   = 200
)

func TestSupplier(t *testing.T) {
	var supplier Supplier[string] = func() string {
		return "Hello"
	}
	assert.Equal(t, "Hello", supplier())
	//
	var fun func() string = supplier
	assert.Equal(t, "Hello", fun())
}

//===

func TestNewThreadLocal_Single(t *testing.T) {
	tls := NewThreadLocal[string]()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2 := NewThreadLocal[int]()
	assert.Equal(t, "Hello", tls.Get())
	tls2.Set(22)
	assert.Equal(t, 22, tls2.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	task := GoWait(func(token CancelToken) {
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, 0, tls2.Get())
	})
	task.Get()
}

func TestNewThreadLocal_Multi(t *testing.T) {
	tls := NewThreadLocal[string]()
	tls2 := NewThreadLocal[int]()
	tls.Set("Hello")
	tls2.Set(22)
	assert.Equal(t, 22, tls2.Get())
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	task := GoWait(func(token CancelToken) {
		assert.Equal(t, "", tls.Get())
		assert.Equal(t, 0, tls2.Get())
	})
	task.Get()
}

func TestNewThreadLocal_Concurrency(t *testing.T) {
	tls := NewThreadLocal[uint64]()
	tls2 := NewThreadLocal[uint64]()
	//
	tls2.Set(33)
	assert.Equal(t, uint64(33), tls2.Get())
	//
	wg := &sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		assert.Equal(t, uint64(0), tls.Get())
		assert.Equal(t, uint64(33), tls2.Get())
		Go(func() {
			assert.Equal(t, uint64(0), tls.Get())
			assert.Equal(t, uint64(0), tls2.Get())
			v := rand.Uint64()
			v2 := rand.Uint64()
			for j := 0; j < loopTimes; j++ {
				tls.Set(v)
				tmp := tls.Get()
				assert.Equal(t, v, tmp)
				//
				tls2.Set(v2)
				tmp2 := tls2.Get()
				assert.Equal(t, v2, tmp2)
			}
			wg.Done()
		})
	}
	wg.Wait()
	//
	task := GoWait(func(token CancelToken) {
		assert.Equal(t, uint64(0), tls.Get())
		assert.Equal(t, uint64(0), tls2.Get())
	})
	task.Get()
}

//===

func TestNewThreadLocalWithInitial_Single(t *testing.T) {
	tls := NewThreadLocalWithInitial[string](func() string {
		return "Hello"
	})
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2 := NewThreadLocalWithInitial[int](func() int {
		return 22
	})
	assert.Equal(t, "Hello", tls.Get())
	assert.Equal(t, 22, tls2.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	task := GoWait(func(token CancelToken) {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 22, tls2.Get())
	})
	task.Get()
}

func TestNewThreadLocalWithInitial_Multi(t *testing.T) {
	tls := NewThreadLocalWithInitial[string](func() string {
		return "Hello"
	})
	tls2 := NewThreadLocalWithInitial[int](func() int {
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
	task := GoWait(func(token CancelToken) {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 22, tls2.Get())
	})
	task.Get()
}

func TestNewThreadLocalWithInitial_Concurrency(t *testing.T) {
	tls := NewThreadLocalWithInitial[any](func() any {
		return "Hello"
	})
	tls2 := NewThreadLocalWithInitial[uint64](func() uint64 {
		return uint64(22)
	})
	//
	tls2.Set(33)
	assert.Equal(t, uint64(33), tls2.Get())
	//
	wg := &sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, uint64(33), tls2.Get())
		Go(func() {
			assert.Equal(t, "Hello", tls.Get())
			assert.Equal(t, uint64(22), tls2.Get())
			v := rand.Uint64()
			v2 := rand.Uint64()
			for j := 0; j < loopTimes; j++ {
				tls.Set(v)
				tmp := tls.Get()
				assert.Equal(t, v, tmp.(uint64))
				//
				tls2.Set(v2)
				tmp2 := tls2.Get()
				assert.Equal(t, v2, tmp2)
			}
			wg.Done()
		})
	}
	wg.Wait()
	//
	task := GoWait(func(token CancelToken) {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, uint64(22), tls2.Get())
	})
	task.Get()
}

//===

func TestNewInheritableThreadLocal_Single(t *testing.T) {
	tls := NewInheritableThreadLocal[string]()
	tls.Set("Hello")
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2 := NewInheritableThreadLocal[int]()
	assert.Equal(t, "Hello", tls.Get())
	tls2.Set(22)
	assert.Equal(t, 22, tls2.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	task := GoWait(func(token CancelToken) {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 33, tls2.Get())
	})
	task.Get()
}

func TestNewInheritableThreadLocal_Multi(t *testing.T) {
	tls := NewInheritableThreadLocal[string]()
	tls2 := NewInheritableThreadLocal[int]()
	tls.Set("Hello")
	tls2.Set(22)
	assert.Equal(t, 22, tls2.Get())
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	task := GoWait(func(token CancelToken) {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 33, tls2.Get())
	})
	task.Get()
}

func TestNewInheritableThreadLocal_Concurrency(t *testing.T) {
	tls := NewInheritableThreadLocal[uint64]()
	tls2 := NewInheritableThreadLocal[uint64]()
	//
	tls2.Set(33)
	assert.Equal(t, uint64(33), tls2.Get())
	//
	wg := &sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		assert.Equal(t, uint64(0), tls.Get())
		assert.Equal(t, uint64(33), tls2.Get())
		Go(func() {
			assert.Equal(t, uint64(0), tls.Get())
			assert.Equal(t, uint64(33), tls2.Get())
			v := rand.Uint64()
			v2 := rand.Uint64()
			for j := 0; j < loopTimes; j++ {
				tls.Set(v)
				tmp := tls.Get()
				assert.Equal(t, v, tmp)
				//
				tls2.Set(v2)
				tmp2 := tls2.Get()
				assert.Equal(t, v2, tmp2)
			}
			wg.Done()
		})
	}
	wg.Wait()
	//
	task := GoWait(func(token CancelToken) {
		assert.Equal(t, uint64(0), tls.Get())
		assert.Equal(t, uint64(33), tls2.Get())
	})
	task.Get()
}

//===

func TestNewInheritableThreadLocalWithInitial_Single(t *testing.T) {
	tls := NewInheritableThreadLocalWithInitial[string](func() string {
		return "Hello"
	})
	assert.Equal(t, "Hello", tls.Get())
	//
	tls2 := NewInheritableThreadLocalWithInitial[int](func() int {
		return 22
	})
	assert.Equal(t, "Hello", tls.Get())
	assert.Equal(t, 22, tls2.Get())
	//
	tls2.Set(33)
	assert.Equal(t, 33, tls2.Get())
	//
	task := GoWait(func(token CancelToken) {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 33, tls2.Get())
	})
	task.Get()
}

func TestNewInheritableThreadLocalWithInitial_Multi(t *testing.T) {
	tls := NewInheritableThreadLocalWithInitial[string](func() string {
		return "Hello"
	})
	tls2 := NewInheritableThreadLocalWithInitial[int](func() int {
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
	task := GoWait(func(token CancelToken) {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, 33, tls2.Get())
	})
	task.Get()
}

func TestNewInheritableThreadLocalWithInitial_Concurrency(t *testing.T) {
	tls := NewInheritableThreadLocalWithInitial[any](func() any {
		return "Hello"
	})
	tls2 := NewInheritableThreadLocalWithInitial[uint64](func() uint64 {
		return uint64(22)
	})
	//
	tls2.Set(33)
	assert.Equal(t, uint64(33), tls2.Get())
	//
	wg := &sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, uint64(33), tls2.Get())
		Go(func() {
			assert.Equal(t, "Hello", tls.Get())
			assert.Equal(t, uint64(33), tls2.Get())
			v := rand.Uint64()
			v2 := rand.Uint64()
			for j := 0; j < loopTimes; j++ {
				tls.Set(v)
				tmp := tls.Get()
				assert.Equal(t, v, tmp.(uint64))
				//
				tls2.Set(v2)
				tmp2 := tls2.Get()
				assert.Equal(t, v2, tmp2)
			}
			wg.Done()
		})
	}
	wg.Wait()
	//
	task := GoWait(func(token CancelToken) {
		assert.Equal(t, "Hello", tls.Get())
		assert.Equal(t, uint64(33), tls2.Get())
	})
	task.Get()
}

//===

// BenchmarkThreadLocal-8                          13636471                94.17 ns/op            7 B/op          0 allocs/op
func BenchmarkThreadLocal(b *testing.B) {
	tlsCount := 100
	tlsSlice := make([]ThreadLocal[int], tlsCount)
	for i := 0; i < tlsCount; i++ {
		tlsSlice[i] = NewThreadLocal[int]()
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % tlsCount
		tls := tlsSlice[index]
		initValue := tls.Get()
		if initValue != 0 {
			b.Fail()
		}
		tls.Set(i)
		if tls.Get() != i {
			b.Fail()
		}
		tls.Remove()
	}
}

// BenchmarkThreadLocalWithInitial-8               13674153                86.76 ns/op            7 B/op          0 allocs/op
func BenchmarkThreadLocalWithInitial(b *testing.B) {
	tlsCount := 100
	tlsSlice := make([]ThreadLocal[int], tlsCount)
	for i := 0; i < tlsCount; i++ {
		index := i
		tlsSlice[i] = NewThreadLocalWithInitial[int](func() int {
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

// BenchmarkInheritableThreadLocal-8               13917819                84.27 ns/op            7 B/op          0 allocs/op
func BenchmarkInheritableThreadLocal(b *testing.B) {
	tlsCount := 100
	tlsSlice := make([]ThreadLocal[int], tlsCount)
	for i := 0; i < tlsCount; i++ {
		tlsSlice[i] = NewInheritableThreadLocal[int]()
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % tlsCount
		tls := tlsSlice[index]
		initValue := tls.Get()
		if initValue != 0 {
			b.Fail()
		}
		tls.Set(i)
		if tls.Get() != i {
			b.Fail()
		}
		tls.Remove()
	}
}

// BenchmarkInheritableThreadLocalWithInitial-8    13483130                90.03 ns/op            7 B/op          0 allocs/op
func BenchmarkInheritableThreadLocalWithInitial(b *testing.B) {
	tlsCount := 100
	tlsSlice := make([]ThreadLocal[int], tlsCount)
	for i := 0; i < tlsCount; i++ {
		index := i
		tlsSlice[i] = NewInheritableThreadLocalWithInitial[int](func() int {
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

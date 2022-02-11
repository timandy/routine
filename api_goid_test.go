package routine

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGoid(t *testing.T) {
	assert.NotEqual(t, 0, Goid())
}

func TestAllGoids(t *testing.T) {
	const num = 10
	for i := 0; i < num; i++ {
		go func() {
			time.Sleep(time.Second)
		}()
	}
	time.Sleep(time.Millisecond)
	assert.NotEmpty(t, AllGoids())
}

func TestForeachGoid(t *testing.T) {
	const num = 10
	for i := 0; i < num; i++ {
		go func() {
			time.Sleep(time.Second)
		}()
	}
	time.Sleep(time.Millisecond)
	//
	cnt := 0
	ForeachGoid(func(goid int64) {
		assert.NotEqual(t, 0, goid)
		cnt++
	})
	assert.GreaterOrEqual(t, cnt, num)
}

// BenchmarkGoid-4                279375390             4.244 ns/op               0 B/op               0 allocs/op
func BenchmarkGoid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Goid()
	}
}

// BenchmarkAllGoids-4               469323              2132 ns/op             640 B/op               1 allocs/op
func BenchmarkAllGoids(b *testing.B) {
	const num = 16
	for i := 0; i < num; i++ {
		go func() {
			time.Sleep(time.Minute)
		}()
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AllGoids()
	}
}

// BenchmarkForeachGoid-4            677912              1943 ns/op               0 B/op               0 allocs/op
func BenchmarkForeachGoid(b *testing.B) {
	const num = 16
	for i := 0; i < num; i++ {
		go func() {
			time.Sleep(time.Minute)
		}()
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ForeachGoid(func(goid int64) {})
	}
}

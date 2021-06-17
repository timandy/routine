package routine

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGoid(t *testing.T) {
	t.Log(Goid())
}

func TestAllGoid(t *testing.T) {
	const num = 10
	for i := 0; i < num; i++ {
		go func() {
			time.Sleep(time.Second)
		}()
	}
	time.Sleep(time.Millisecond)

	ids := AllGoids()
	t.Log("all gids: ", len(ids), ids)
}

func TestGoStorage(t *testing.T) {
	var variable = "hello world"
	stg := NewLocalStorage()
	stg.Set(variable)
	Go(func() {
		v := stg.Get()
		assert.True(t, v != nil && v.(string) == variable)
	})
	time.Sleep(time.Millisecond)
	stg.Clear()
}

// BenchmarkGoid-12    	278801190	         4.586 ns/op	       0 B/op	       0 allocs/op
func BenchmarkGoid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Goid()
	}
}

// BenchmarkAllGoid-12    	 5949680	       228.3 ns/op	     896 B/op	       1 allocs/op
func BenchmarkAllGoid(b *testing.B) {
	const num = 16
	for i := 0; i < num; i++ {
		go func() {
			time.Sleep(time.Second)
		}()
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AllGoids()
	}
}

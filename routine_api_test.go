package routine

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNewLocalStorage(t *testing.T) {
	s := NewLocalStorage()
	s.Set("hello")
	assert.Equal(t, "hello", s.Get())
	//
	s2 := NewLocalStorage()
	assert.Equal(t, "hello", s.Get())
	s2.Set(22)
	assert.Equal(t, 22, s2.Get())
}

func TestMultiStorage(t *testing.T) {
	s := NewLocalStorage()
	s2 := NewLocalStorage()
	s.Set("hello")
	s2.Set(22)
	assert.Equal(t, 22, s2.Get())
	assert.Equal(t, "hello", s.Get())
}

func TestBackupContext(t *testing.T) {
	s := NewLocalStorage()
	ic := BackupContext()

	waiter := &sync.WaitGroup{}
	waiter.Add(1)
	go func() {
		s.Set("hello")
		assert.Equal(t, "hello", s.Get())
		icLocalBackup := BackupContext()
		//
		InheritContext(ic)
		assert.Nil(t, s.Get())
		//
		InheritContext(icLocalBackup)
		assert.Equal(t, "hello", s.Get())
		//
		waiter.Done()
	}()
	waiter.Wait()
}

func TestGoid(t *testing.T) {
	assert.NotEqual(t, 0, Goid())
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
	waiter := &sync.WaitGroup{}
	waiter.Add(1)
	variable := "hello world"
	stg := NewLocalStorage()
	stg.Set(variable)
	Go(func() {
		v := stg.Get()
		assert.Equal(t, variable, v.(string))
		waiter.Done()
	})
	waiter.Wait()
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

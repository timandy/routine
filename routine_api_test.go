package routine

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNewThreadLocal(t *testing.T) {
	threadLocal := NewThreadLocal()
	threadLocal.Set("hello")
	assert.Equal(t, "hello", threadLocal.Get())
	//
	threadLocal2 := NewThreadLocal()
	assert.Equal(t, "hello", threadLocal.Get())
	threadLocal2.Set(22)
	assert.Equal(t, 22, threadLocal2.Get())
}

func TestMultiThreadLocal(t *testing.T) {
	threadLocal := NewThreadLocal()
	threadLocal2 := NewThreadLocal()
	threadLocal.Set("hello")
	threadLocal2.Set(22)
	assert.Equal(t, 22, threadLocal2.Get())
	assert.Equal(t, "hello", threadLocal.Get())
}

func TestBackupContext(t *testing.T) {
	threadLocal := NewThreadLocal()
	ic := BackupContext()

	waiter := &sync.WaitGroup{}
	waiter.Add(1)
	go func() {
		threadLocal.Set("hello")
		assert.Equal(t, "hello", threadLocal.Get())
		icLocalBackup := BackupContext()
		//
		RestoreContext(ic)
		assert.Nil(t, threadLocal.Get())
		//
		RestoreContext(icLocalBackup)
		assert.Equal(t, "hello", threadLocal.Get())
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

func TestGoThreadLocal(t *testing.T) {
	waiter := &sync.WaitGroup{}
	waiter.Add(1)
	variable := "hello world"
	threadLocal := NewThreadLocal()
	threadLocal.Set(variable)
	Go(func() {
		v := threadLocal.Get()
		assert.Equal(t, variable, v.(string))
		waiter.Done()
	})
	waiter.Wait()
}

func TestClear(t *testing.T) {
	threadLocal := NewThreadLocal()
	Clear()
	assert.Nil(t, threadLocal.Get())
	threadLocal.Set(1)
	assert.Equal(t, 1, threadLocal.Get())
	Clear()
	assert.Nil(t, threadLocal.Get())
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

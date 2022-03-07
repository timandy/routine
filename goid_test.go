package routine

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestGetGoidByNative(t *testing.T) {
	sid := getGoidByStack()
	assert.NotEqual(t, 0, sid)
	//
	for i := 0; i < 100; i++ {
		nid, success := getGoidByNative()
		assert.True(t, success)
		assert.Equal(t, sid, nid)
	}
	//
	wg := &sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			sid2 := getGoidByStack()
			assert.NotEqual(t, 0, sid2)
			//
			nid2, success2 := getGoidByNative()
			assert.True(t, success2)
			assert.Equal(t, sid2, nid2)
			assert.NotEqual(t, sid, nid2)
			//
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestGetGoidByStack(t *testing.T) {
	nid, success := getGoidByNative()
	assert.True(t, success)
	assert.NotEqual(t, 0, nid)
	//
	for i := 0; i < 100; i++ {
		sid := getGoidByStack()
		assert.Equal(t, nid, sid)
	}
	//
	wg := &sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			nid2, success2 := getGoidByNative()
			assert.True(t, success2)
			assert.NotEqual(t, 0, nid2)
			//
			sid2 := getGoidByStack()
			assert.Equal(t, nid2, sid2)
			assert.NotEqual(t, nid, sid2)
			//
			wg.Done()
		}()
	}
	wg.Wait()
}

//===

// BenchmarkGetGoidByNative-4                     363983250                3.270 ns/op            0 B/op          0 allocs/op
func BenchmarkGetGoidByNative(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = getGoidByNative()
	}
}

// BenchmarkGetGoidByStack-4                         363614                 3030 ns/op            0 B/op          0 allocs/op
func BenchmarkGetGoidByStack(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getGoidByStack()
	}
}

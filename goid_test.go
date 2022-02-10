package routine

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
	"unsafe"
)

func TestFindGoidPointer(t *testing.T) {
	assert.Nil(t, findGoidPointer(nil))
	//goland:noinspection GoVetUnsafePointer
	assert.Nil(t, findGoidPointer(unsafe.Pointer(uintptr(0))))
}

func TestFindNextGoid(t *testing.T) {
	stack := "goroutine 6 [running]:\n...\ngoroutine 33 [running]..."
	goid, next := findNextGoid([]byte(stack), 0)
	assert.Equal(t, int64(6), goid)
	assert.Equal(t, 12, next)
	//
	goid, next = findNextGoid([]byte(stack), next)
	assert.Equal(t, int64(33), goid)
	assert.Equal(t, 40, next)
}

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

func TestGetAllGoidByNative(t *testing.T) {
	goid, success := getGoidByNative()
	assert.True(t, success)
	assert.NotEqual(t, 0, goid)
	//
	nids, success := getAllGoidByNative()
	if !success {
		return
	}
	//
	for _, nid := range nids {
		if nid == goid {
			return
		}
	}
	assert.Fail(t, "nids must contains current goid")
}

func TestGetAllGoidByStack(t *testing.T) {
	goid, success := getGoidByNative()
	assert.True(t, success)
	assert.NotEqual(t, 0, goid)
	//
	sids := getAllGoidByStack()
	//
	for _, sid := range sids {
		if sid == goid {
			return
		}
	}
	assert.Fail(t, "sids must contains current goid")
}

//===

// BenchmarkGetGoidByNative-4                      279528649                4.303 ns/op           0 B/op          0 allocs/op
func BenchmarkGetGoidByNative(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getGoidByNative()
	}
}

// BenchmarkGetGoidByStack-4                         255331              4898 ns/op              64 B/op          1 allocs/op
func BenchmarkGetGoidByStack(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getGoidByStack()
	}
}

// BenchmarkGetAllGoidByNative-4                        100          16360789 ns/op         1056768 B/op          1 allocs/op
func BenchmarkGetAllGoidByNative(b *testing.B) {
	const routineNum = 65536
	for i := 0; i < routineNum; i++ {
		go func() {
			time.Sleep(time.Minute)
		}()
	}
	time.Sleep(time.Millisecond * 100)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = getAllGoidByNative()
	}
}

// BenchmarkGetAllGoidByStack-4                           1        4925942800 ns/op        142974848 B/op        34 allocs/op
func BenchmarkGetAllGoidByStack(b *testing.B) {
	const routineNum = 65536
	for i := 0; i < routineNum; i++ {
		go func() {
			time.Sleep(time.Minute)
		}()
	}
	time.Sleep(time.Millisecond * 100)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = getAllGoidByStack()
	}
}

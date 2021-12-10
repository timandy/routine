package routine

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"time"
)

func TestGetGoid(t *testing.T) {
	var id int64
	for i := 0; i < 100; i++ {
		nid, _ := getGoidByNative()
		sid := getGoidByStack()
		assert.Equal(t, sid, nid)

		if id == 0 {
			id = sid
		} else {
			assert.Equal(t, sid, id)
		}
	}
	t.Log(getGoidByStack())
	t.Log(getGoidByNative())
}

func TestPC(t *testing.T) {
	const num = 10

	for i := 0; i < num; i++ {
		go func() {
			time.Sleep(time.Minute)
		}()
	}

	time.Sleep(time.Millisecond)
	// can only load pc info, cannot load goroutine context
	res := make([]runtime.StackRecord, 1024)
	n, ok := runtime.GoroutineProfile(res)
	t.Log(n, ok)

	// 获取全部协程上下文快照
	buf := make([]byte, 1<<20)
	n = runtime.Stack(buf, true)
	t.Log("all: \n", string(buf[:n]))
}

// BenchmarkGetGoidByNative-12    	409483110	         2.530 ns/op	       0 B/op	       0 allocs/op
func BenchmarkGetGoidByNative(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getGoidByNative()
	}
}

// BenchmarkGetGoidByStack-12    	  391472	      2996 ns/op	       0 B/op	       0 allocs/op
func BenchmarkGetGoidByStack(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getGoidByStack()
	}
}

// When routineNum = 16
// BenchmarkGetAllGoidByNative-12    	 6226236	       242.9 ns/op	     896 B/op	       1 allocs/op
// When routineNum = 64
// BenchmarkGetAllGoidByNative-12    	 2064740	       740.3 ns/op	    3072 B/op	       1 allocs/op
// When routineNum = 256
// BenchmarkGetAllGoidByNative-12    	  402352	      3163 ns/op	    9472 B/op	       1 allocs/op
// When routineNum = 1024
// BenchmarkGetAllGoidByNative-12    	   84907	     18508 ns/op	   40960 B/op	       1 allocs/op
// When routineNum = 8192
// BenchmarkGetAllGoidByNative-12    	    7520	    253149 ns/op	  270336 B/op	       1 allocs/op
// When routineNum = 65536
// BenchmarkGetAllGoidByNative-12    	     457	   4195487 ns/op	 1581056 B/op	       1 allocs/op
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

// When routineNum = 16
// BenchmarkGetAllGoidByStack-12    	    7779	    147593 ns/op	   57793 B/op	       2 allocs/op
// When routineNum = 64
// BenchmarkGetAllGoidByStack-12    	    2754	    463350 ns/op	  206595 B/op	       2 allocs/op
// When routineNum = 256
// BenchmarkGetAllGoidByStack-12    	     843	   1858839 ns/op	  801156 B/op	       2 allocs/op
// When routineNum = 1024
// BenchmarkGetAllGoidByStack-12    	     254	   6923118 ns/op	 3181195 B/op	       2 allocs/op
// When routineNum = 8192
// BenchmarkGetAllGoidByStack-12    	      43	  37269949 ns/op	16924678 B/op	       2 allocs/op
// When routineNum = 65536
// BenchmarkGetAllGoidByStack-12    	       6	 316648415 ns/op	135282688 B/op	       2 allocs/op
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

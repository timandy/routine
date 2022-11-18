package routine

import (
	"reflect"
	"runtime"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestG_Goid(t *testing.T) {
	runTest(t, func() {
		gp := getg()
		runtime.GC()
		assert.Equal(t, curGoroutineID(), gp.goid)
	})
}

func TestG_Paniconfault(t *testing.T) {
	runTest(t, func() {
		gp := getg()
		runtime.GC()
		//read-1
		assert.False(t, setPanicOnFault(false))
		assert.False(t, gp.getPanicOnFault())
		//read-2
		setPanicOnFault(true)
		assert.True(t, gp.getPanicOnFault())
		//write-1
		gp.setPanicOnFault(false)
		assert.False(t, setPanicOnFault(false))
		//write-2
		gp.setPanicOnFault(true)
		assert.True(t, setPanicOnFault(true))
		//write-read-1
		gp.setPanicOnFault(false)
		assert.False(t, gp.getPanicOnFault())
		//write-read-2
		gp.setPanicOnFault(true)
		assert.True(t, gp.getPanicOnFault())
		//restore
		gp.setPanicOnFault(false)
	})
}

func TestG_Gopc(t *testing.T) {
	runTest(t, func() {
		gp := getg()
		runtime.GC()
		assert.Greater(t, int64(*gp.gopc), int64(0))
	})
}

func TestG_ProfLabel(t *testing.T) {
	runTest(t, func() {
		ptr := unsafe.Pointer(&struct{}{})
		null := unsafe.Pointer(nil)
		assert.NotEqual(t, ptr, null)
		//
		gp := getg()
		runtime.GC()
		//read-1
		assert.Equal(t, null, getProfLabel())
		assert.Equal(t, null, gp.getLabels())
		//read-2
		setProfLabel(ptr)
		assert.Equal(t, ptr, gp.getLabels())
		//write-1
		gp.setLabels(nil)
		assert.Equal(t, null, getProfLabel())
		//write-2
		gp.setLabels(ptr)
		assert.Equal(t, ptr, getProfLabel())
		//write-read-1
		gp.setLabels(nil)
		assert.Equal(t, null, gp.getLabels())
		//write-read-2
		gp.setLabels(ptr)
		assert.Equal(t, ptr, gp.getLabels())
		//restore
		gp.setLabels(null)
	})
}

func TestOffset(t *testing.T) {
	runTest(t, func() {
		assert.Panics(t, func() {
			gt := reflect.TypeOf(0)
			offset(gt, "hello")
		})
		assert.PanicsWithValue(t, "No such field 'hello' of struct 'runtime.g'.", func() {
			gt := getgt()
			offset(gt, "hello")
		})
	})
}

// curGoroutineID parse the current g's goid from caller stack.
//
//go:linkname curGoroutineID net/http.http2curGoroutineID
func curGoroutineID() int64

// setPanicOnFault controls the runtime's behavior when a program faults at an unexpected (non-nil) address.
//
//go:linkname setPanicOnFault runtime/debug.setPanicOnFault
func setPanicOnFault(new bool) (old bool)

// getProfLabel get current g's labels which will be inherited by new goroutine.
//
//go:linkname getProfLabel runtime/pprof.runtime_getProfLabel
func getProfLabel() unsafe.Pointer

// setProfLabel set current g's labels which will be inherited by new goroutine.
//
//go:linkname setProfLabel runtime/pprof.runtime_setProfLabel
func setProfLabel(labels unsafe.Pointer)

//===

// BenchmarkGohack-8                              186637413                5.734 ns/op            0 B/op          0 allocs/op
func BenchmarkGohack(b *testing.B) {
	_ = getg()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gp := getg()
		_ = gp.goid
		_ = gp.gopc
		_ = gp.getLabels()
		_ = gp.getPanicOnFault()
		gp.setLabels(nil)
		gp.setPanicOnFault(false)
	}
}

package routine

import (
	_ "github.com/timandy/routine/g"
	_ "net/http"
	_ "runtime/pprof"
	"unsafe"
)

// getg returns the pointer to the current runtime.g.
//go:linkname getg github.com/timandy/routine/g.getg
func getg() unsafe.Pointer

// curGoroutineID parse the current g's goid from caller stack.
//go:linkname curGoroutineID net/http.http2curGoroutineID
func curGoroutineID() uint64

// getProfLabel get current g's labels which will be inherited by new goroutine.
//go:linkname getProfLabel runtime/pprof.runtime_getProfLabel
func getProfLabel() unsafe.Pointer

// setProfLabel set current g's labels which will be inherited by new goroutine.
//go:linkname setProfLabel runtime/pprof.runtime_setProfLabel
func setProfLabel(labels unsafe.Pointer)

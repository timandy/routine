//go:build go1.12
// +build go1.12

package routine

import "unsafe"

const (
	gDead = 6
)

//go:linkname runtimeG runtime.g
type runtimeG struct {
}

//go:linkname runtimeMutex runtime.mutex
type runtimeMutex struct {
}

//go:linkname runtimeAllgs runtime.allgs
var runtimeAllgs []*runtimeG

//go:linkname runtimeAllglock runtime.allglock
var runtimeAllglock runtimeMutex

//go:linkname runtimeLock runtime.lock
func runtimeLock(l *runtimeMutex)

//go:linkname runtimeUnlock runtime.unlock
func runtimeUnlock(l *runtimeMutex)

//go:linkname runtimeReadgstatus runtime.readgstatus
func runtimeReadgstatus(g *runtimeG) uint32

//go:linkname runtimeIsSystemGoroutine runtime.isSystemGoroutine
func runtimeIsSystemGoroutine(gp *runtimeG, fixed bool) bool

// getAllGoidByNative retrieve all goid by runtime.allgs
func getAllGoidByNative() ([]int64, bool) {
	defer func() {
		recover()
	}()
	runtimeLock(&runtimeAllglock)
	defer runtimeUnlock(&runtimeAllglock)
	allgs := runtimeAllgs
	goids := make([]int64, 0, len(allgs))
	for _, gp := range allgs {
		if runtimeReadgstatus(gp) == gDead || runtimeIsSystemGoroutine(gp, false) {
			continue
		}
		goid := findGoidPointer(unsafe.Pointer(gp))
		if goid == nil {
			continue
		}
		goids = append(goids, *goid)
	}
	return goids, true
}

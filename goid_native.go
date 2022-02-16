//go:build go1.13
// +build go1.13

package routine

import "unsafe"

const (
	gDead    = 6
	goidSize = 1024
)

//go:linkname g runtime.g
type g struct {
}

//go:linkname mutex runtime.mutex
type mutex struct {
}

//go:linkname allgs runtime.allgs
var allgs []*g

//go:linkname allglock runtime.allglock
var allglock mutex

//go:linkname lock runtime.lock
func lock(l *mutex)

//go:linkname unlock runtime.unlock
func unlock(l *mutex)

//go:linkname readgstatus runtime.readgstatus
func readgstatus(gp *g) uint32

//go:linkname isSystemGoroutine runtime.isSystemGoroutine
func isSystemGoroutine(gp *g, fixed bool) bool

// atomicAllG return allgs safely under the protection of allglock.
// New Gs appended during the race can be missed.
func atomicAllG() []*g {
	lock(&allglock)
	defer unlock(&allglock)
	return allgs
}

// getAllGoidByNative retrieve all goid through native.
// Addition of new Gs during execution, which may be missed.
func getAllGoidByNative() ([]int64, bool) {
	if !support() {
		return nil, false
	}
	allg := atomicAllG()
	goids := make([]int64, 0, goidSize)
	for _, gp := range allg {
		if readgstatus(gp) == gDead || isSystemGoroutine(gp, false) {
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

// foreachGoidByNative run a func for each goroutine's goid through native.
// Addition of new Gs during execution, which may be missed.
func foreachGoidByNative(fun func(goid int64)) bool {
	if !support() {
		return false
	}
	allg := atomicAllG()
	for _, gp := range allg {
		if readgstatus(gp) == gDead || isSystemGoroutine(gp, false) {
			continue
		}
		goid := findGoidPointer(unsafe.Pointer(gp))
		if goid == nil {
			continue
		}
		fun(*goid)
	}
	return true
}

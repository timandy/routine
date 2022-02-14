//go:build go1.13
// +build go1.13

package routine

import "unsafe"

const (
	gDead = 6
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

// getAllGoidByNative retrieve all goid through runtime.allgs
func getAllGoidByNative() ([]int64, bool) {
	if !support() {
		return nil, false
	}
	lock(&allglock)
	defer unlock(&allglock)
	goids := make([]int64, 0, len(allgs))
	for _, gp := range allgs {
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

// foreachGoidByNative run a func for each goroutine's goid through runtime.allgs
func foreachGoidByNative(fun func(goid int64)) bool {
	if !support() {
		return false
	}
	lock(&allglock)
	defer unlock(&allglock)
	for _, gp := range allgs {
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

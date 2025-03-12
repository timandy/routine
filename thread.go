//go:build !routinex

package routine

import (
	"runtime"
	"unsafe"
)

const threadMagic = uint64('r')<<48 |
	uint64('o')<<40 |
	uint64('u')<<32 |
	uint64('t')<<24 |
	uint64('i')<<16 |
	uint64('n')<<8 |
	uint64('e')

type thread struct {
	labels                  labelMap //pprof
	magic                   uint64   //mark
	id                      uint64   //goid
	threadLocals            *threadLocalMap
	inheritableThreadLocals *threadLocalMap
}

// finalize reset thread's memory.
func (t *thread) finalize() {
	t.labels = nil
	t.magic = 0
	t.id = 0
	t.threadLocals = nil
	t.inheritableThreadLocals = nil
}

// currentThread returns a pointer to the currently executing goroutine's thread struct.
//
//go:norace
//go:nocheckptr
func currentThread(create bool) *thread {
	gp := getg()
	goid := gp.goid()
	label := gp.getLabels()
	//nothing inherited
	if label == nil {
		if create {
			newt := &thread{labels: nil, magic: threadMagic, id: goid}
			runtime.SetFinalizer(newt, (*thread).finalize)
			gp.setLabels(unsafe.Pointer(newt))
			return newt
		}
		return nil
	}
	//inherited map then create
	t, magic, id := extractThread(gp, label)
	if magic != threadMagic {
		if create {
			mp := *(*labelMap)(label)
			newt := &thread{labels: mp, magic: threadMagic, id: goid}
			runtime.SetFinalizer(newt, (*thread).finalize)
			gp.setLabels(unsafe.Pointer(newt))
			return newt
		}
		return nil
	}
	//inherited thread then recreate
	if id != goid {
		if create || t.labels != nil {
			newt := &thread{labels: t.labels, magic: threadMagic, id: goid}
			runtime.SetFinalizer(newt, (*thread).finalize)
			gp.setLabels(unsafe.Pointer(newt))
			return newt
		}
		gp.setLabels(nil)
		return nil
	}
	//all is ok
	return t
}

// extractThread extract thread from unsafe.Pointer and catch fault error.
//
//go:norace
//go:nocheckptr
func extractThread(gp *g, label unsafe.Pointer) (t *thread, magic uint64, id uint64) {
	old := gp.setPanicOnFault(true)
	defer func() {
		gp.setPanicOnFault(old)
		recover() //nolint:errcheck
	}()
	t = (*thread)(label)
	return t, t.magic, t.id
}

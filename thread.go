package routine

import (
	"unsafe"
)

type thread struct {
	id                      int64
	threadLocals            *threadLocalMap
	inheritableThreadLocals *threadLocalMap
}

func currentThread() *thread {
	goid := Goid()
	label := getProfLabel()
	//nothing inherited
	if label == nil {
		newt := &thread{
			id: goid,
		}
		setProfLabel(unsafe.Pointer(newt))
		return newt
	}
	//inherited need recreate
	t := (*thread)(label)
	if t.id != goid {
		newt := &thread{
			id:                      goid,
			inheritableThreadLocals: createInheritedMap(t.inheritableThreadLocals),
		}
		setProfLabel(unsafe.Pointer(newt))
		return newt
	}
	//all is ok
	return t
}

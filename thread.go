package routine

import "unsafe"

type thread struct {
	id                      int64
	threadLocals            *threadLocalMap
	inheritableThreadLocals *threadLocalMap
}

func currentThread(create bool) *thread {
	goid := Goid()
	label := getProfLabel()
	//nothing inherited
	if label == nil {
		if create {
			newt := &thread{id: goid}
			setProfLabel(unsafe.Pointer(newt))
			return newt
		}
		return nil
	}
	//inherited need recreate
	t := (*thread)(label)
	if t.id != goid {
		if create {
			newt := &thread{id: goid}
			setProfLabel(unsafe.Pointer(newt))
			return newt
		}
		setProfLabel(nil)
		return nil
	}
	//all is ok
	return t
}

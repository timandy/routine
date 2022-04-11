package routine

import "unsafe"

const threadMagic = int64('r')<<48 |
	int64('o')<<40 |
	int64('u')<<32 |
	int64('t')<<24 |
	int64('i')<<16 |
	int64('n')<<8 |
	int64('e')

type thread struct {
	magic                   int64
	id                      int64
	threadLocals            *threadLocalMap
	inheritableThreadLocals *threadLocalMap
}

func currentThread(create bool) *thread {
	gp := getg()
	goid := gp.goid
	label := gp.getLabels()
	//nothing inherited
	if label == nil {
		if create {
			newt := &thread{magic: threadMagic, id: goid}
			gp.setLabels(unsafe.Pointer(newt))
			return newt
		}
		return nil
	}
	//inherited need recreate
	t := (*thread)(label)
	if t.id != goid || t.magic != threadMagic {
		if create {
			newt := &thread{magic: threadMagic, id: goid}
			gp.setLabels(unsafe.Pointer(newt))
			return newt
		}
		gp.setLabels(nil)
		return nil
	}
	//all is ok
	return t
}

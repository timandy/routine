package routine

import "sync/atomic"

var inheritableThreadLocalIndex int32 = -1

func nextInheritableThreadLocalIndex() int {
	index := atomic.AddInt32(&inheritableThreadLocalIndex, 1)
	if index < 0 {
		atomic.AddInt32(&inheritableThreadLocalIndex, -1)
		panic("too many inheritable-thread-local indexed variables")
	}
	return int(index)
}

type inheritableThreadLocal struct {
	index    int
	supplier Supplier
}

func (tls *inheritableThreadLocal) Get() any {
	t := currentThread(true)
	mp := tls.getMap(t)
	if mp != nil {
		v := mp.get(tls.index)
		if v != unset {
			return v
		}
	}
	return tls.setInitialValue(t)
}

func (tls *inheritableThreadLocal) Set(value any) {
	t := currentThread(true)
	mp := tls.getMap(t)
	if mp != nil {
		mp.set(tls.index, value)
	} else {
		tls.createMap(t, value)
	}
}

func (tls *inheritableThreadLocal) Remove() {
	t := currentThread(false)
	if t == nil {
		return
	}
	mp := tls.getMap(t)
	if mp != nil {
		mp.remove(tls.index)
	}
}

func (tls *inheritableThreadLocal) getMap(t *thread) *threadLocalMap {
	return t.inheritableThreadLocals
}

func (tls *inheritableThreadLocal) createMap(t *thread, firstValue any) {
	mp := &threadLocalMap{}
	mp.set(tls.index, firstValue)
	t.inheritableThreadLocals = mp
}

func (tls *inheritableThreadLocal) setInitialValue(t *thread) any {
	value := tls.initialValue()
	mp := tls.getMap(t)
	if mp != nil {
		mp.set(tls.index, value)
	} else {
		tls.createMap(t, value)
	}
	return value
}

func (tls *inheritableThreadLocal) initialValue() any {
	if tls.supplier == nil {
		return nil
	}
	return tls.supplier()
}

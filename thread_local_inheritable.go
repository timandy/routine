package routine

import "sync/atomic"

var inheritableThreadLocalIndex int32 = -1

func nextInheritableThreadLocalId() int {
	index := atomic.AddInt32(&inheritableThreadLocalIndex, 1)
	if index < 0 {
		panic("too many inheritable-thread-local indexed variables")
	}
	return int(index)
}

type inheritableThreadLocal struct {
	id       int
	supplier func() Any
}

func (tls *inheritableThreadLocal) Id() int {
	return tls.id
}

func (tls *inheritableThreadLocal) Get() Any {
	t := currentThread()
	mp := tls.getMap(t)
	if mp != nil {
		v := mp.get(tls)
		if v != unset {
			return v
		}
	}
	return tls.setInitialValue(t)
}

func (tls *inheritableThreadLocal) Set(value Any) {
	t := currentThread()
	mp := tls.getMap(t)
	if mp != nil {
		mp.set(tls, value)
	} else {
		tls.createMap(t, value)
	}
}

func (tls *inheritableThreadLocal) Remove() {
	t := currentThread()
	mp := tls.getMap(t)
	if mp != nil {
		mp.remove(tls)
	}
}

func (tls *inheritableThreadLocal) getMap(t *thread) *threadLocalMap {
	return t.inheritableThreadLocals
}

func (tls *inheritableThreadLocal) createMap(t *thread, firstValue Any) {
	mp := &threadLocalMap{}
	mp.set(tls, firstValue)
	t.inheritableThreadLocals = mp
}

func (tls *inheritableThreadLocal) setInitialValue(t *thread) Any {
	value := tls.initialValue()
	mp := tls.getMap(t)
	if mp != nil {
		mp.set(tls, value)
	} else {
		tls.createMap(t, value)
	}
	return value
}

func (tls *inheritableThreadLocal) initialValue() Any {
	if tls.supplier == nil {
		return nil
	}
	return tls.supplier()
}

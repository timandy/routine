package routine

import "sync/atomic"

var threadLocalIndex int32 = -1

func nextThreadLocalId() int {
	index := atomic.AddInt32(&threadLocalIndex, 1)
	if index < 0 {
		panic("too many thread-local indexed variables")
	}
	return int(index)
}

type threadLocal struct {
	id       int
	supplier func() Any
}

func (tls *threadLocal) Id() int {
	return tls.id
}

func (tls *threadLocal) Get() Any {
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

func (tls *threadLocal) Set(value Any) {
	t := currentThread()
	mp := tls.getMap(t)
	if mp != nil {
		mp.set(tls, value)
	} else {
		tls.createMap(t, value)
	}
}

func (tls *threadLocal) Remove() {
	t := currentThread()
	mp := tls.getMap(t)
	if mp != nil {
		mp.remove(tls)
	}
}

func (tls *threadLocal) getMap(t *thread) *threadLocalMap {
	return t.threadLocals
}

func (tls *threadLocal) createMap(t *thread, firstValue Any) {
	mp := &threadLocalMap{}
	mp.set(tls, firstValue)
	t.threadLocals = mp
}

func (tls *threadLocal) setInitialValue(t *thread) Any {
	value := tls.initialValue()
	mp := tls.getMap(t)
	if mp != nil {
		mp.set(tls, value)
	} else {
		tls.createMap(t, value)
	}
	return value
}

func (tls *threadLocal) initialValue() Any {
	if tls.supplier == nil {
		return nil
	}
	return tls.supplier()
}

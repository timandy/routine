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

type threadLocalImpl struct {
	id       int
	supplier func() Any
}

func (tls *threadLocalImpl) Id() int {
	return tls.id
}

func (tls *threadLocalImpl) Get() Any {
	t := currentThread(true)
	mp := tls.getMap(t)
	if mp != nil {
		e := mp.getEntry(tls)
		if e != nil {
			return e.value
		}
	}
	return tls.setInitialValue(t)
}

func (tls *threadLocalImpl) Set(value Any) {
	t := currentThread(true)
	mp := tls.getMap(t)
	if mp != nil {
		mp.set(tls, value)
	} else {
		tls.createMap(t, value)
	}
}

func (tls *threadLocalImpl) Remove() {
	t := currentThread(false)
	if t == nil {
		return
	}
	mp := tls.getMap(t)
	if mp != nil {
		mp.remove(tls)
	}
}

func (tls *threadLocalImpl) getMap(t *thread) *threadLocalMap {
	return t.threadLocals
}

func (tls *threadLocalImpl) createMap(t *thread, firstValue Any) {
	mp := &threadLocalMap{}
	mp.set(tls, firstValue)
	t.threadLocals = mp
}

func (tls *threadLocalImpl) setInitialValue(t *thread) Any {
	value := tls.initialValue()
	mp := tls.getMap(t)
	if mp != nil {
		mp.set(tls, value)
	} else {
		tls.createMap(t, value)
	}
	return value
}

func (tls *threadLocalImpl) initialValue() Any {
	if tls.supplier == nil {
		return nil
	}
	return tls.supplier()
}

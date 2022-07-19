package routine

import "sync/atomic"

var threadLocalIndex int32 = -1

func nextThreadLocalIndex() int {
	index := atomic.AddInt32(&threadLocalIndex, 1)
	if index < 0 {
		atomic.AddInt32(&threadLocalIndex, -1)
		panic("too many thread-local indexed variables")
	}
	return int(index)
}

type threadLocal struct {
	index    int
	supplier Supplier
}

func (tls *threadLocal) Get() any {
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

func (tls *threadLocal) Set(value any) {
	t := currentThread(true)
	mp := tls.getMap(t)
	if mp != nil {
		mp.set(tls.index, value)
	} else {
		tls.createMap(t, value)
	}
}

func (tls *threadLocal) Remove() {
	t := currentThread(false)
	if t == nil {
		return
	}
	mp := tls.getMap(t)
	if mp != nil {
		mp.remove(tls.index)
	}
}

func (tls *threadLocal) getMap(t *thread) *threadLocalMap {
	return t.threadLocals
}

func (tls *threadLocal) createMap(t *thread, firstValue any) {
	mp := &threadLocalMap{}
	mp.set(tls.index, firstValue)
	t.threadLocals = mp
}

func (tls *threadLocal) setInitialValue(t *thread) any {
	value := tls.initialValue()
	mp := tls.getMap(t)
	if mp != nil {
		mp.set(tls.index, value)
	} else {
		tls.createMap(t, value)
	}
	return value
}

func (tls *threadLocal) initialValue() any {
	if tls.supplier == nil {
		return nil
	}
	return tls.supplier()
}

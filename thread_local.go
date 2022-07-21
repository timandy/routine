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

type threadLocal[T any] struct {
	index    int
	supplier Supplier[T]
}

func (tls *threadLocal[T]) Get() T {
	t := currentThread(true)
	mp := tls.getMap(t)
	if mp != nil {
		v := mp.get(tls.index)
		if v != unset {
			return v.(T)
		}
	}
	return tls.setInitialValue(t)
}

func (tls *threadLocal[T]) Set(value T) {
	t := currentThread(true)
	mp := tls.getMap(t)
	if mp != nil {
		mp.set(tls.index, value)
	} else {
		tls.createMap(t, value)
	}
}

func (tls *threadLocal[T]) Remove() {
	t := currentThread(false)
	if t == nil {
		return
	}
	mp := tls.getMap(t)
	if mp != nil {
		mp.remove(tls.index)
	}
}

func (tls *threadLocal[T]) getMap(t *thread) *threadLocalMap {
	return t.threadLocals
}

func (tls *threadLocal[T]) createMap(t *thread, firstValue T) {
	mp := &threadLocalMap{}
	mp.set(tls.index, firstValue)
	t.threadLocals = mp
}

func (tls *threadLocal[T]) setInitialValue(t *thread) T {
	value := tls.initialValue()
	mp := tls.getMap(t)
	if mp != nil {
		mp.set(tls.index, value)
	} else {
		tls.createMap(t, value)
	}
	return value
}

func (tls *threadLocal[T]) initialValue() T {
	if tls.supplier == nil {
		var defaultValue T
		return defaultValue
	}
	return tls.supplier()
}

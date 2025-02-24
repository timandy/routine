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

type inheritableThreadLocal[T any] struct {
	index    int
	supplier Supplier[T]
}

func (tls *inheritableThreadLocal[T]) Get() T {
	t := currentThread(true)
	mp := tls.getMap(t)
	if mp != nil {
		v := mp.get(tls.index)
		if v != unset {
			return entryValue[T](v)
		}
	}
	return tls.setInitialValue(t)
}

func (tls *inheritableThreadLocal[T]) Set(value T) {
	t := currentThread(true)
	mp := tls.getMap(t)
	if mp != nil {
		mp.set(tls.index, entry(value))
	} else {
		tls.createMap(t, value)
	}
}

func (tls *inheritableThreadLocal[T]) Remove() {
	t := currentThread(false)
	if t == nil {
		return
	}
	mp := tls.getMap(t)
	if mp != nil {
		mp.remove(tls.index)
	}
}

//go:norace
func (tls *inheritableThreadLocal[T]) getMap(t *thread) *threadLocalMap {
	return t.inheritableThreadLocals
}

//go:norace
func (tls *inheritableThreadLocal[T]) createMap(t *thread, firstValue T) {
	mp := &threadLocalMap{}
	mp.set(tls.index, entry(firstValue))
	t.inheritableThreadLocals = mp
}

func (tls *inheritableThreadLocal[T]) setInitialValue(t *thread) T {
	value := tls.initialValue()
	mp := tls.getMap(t)
	if mp != nil {
		mp.set(tls.index, entry(value))
	} else {
		tls.createMap(t, value)
	}
	return value
}

func (tls *inheritableThreadLocal[T]) initialValue() T {
	if tls.supplier == nil {
		var defaultValue T
		return defaultValue
	}
	return tls.supplier()
}

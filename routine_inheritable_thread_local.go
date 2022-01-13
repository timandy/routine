package routine

type inheritableThreadLocalImpl struct {
	id       int
	supplier func() interface{}
}

func (tls *inheritableThreadLocalImpl) Id() int {
	return tls.id
}

func (tls *inheritableThreadLocalImpl) Get() interface{} {
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

func (tls *inheritableThreadLocalImpl) Set(value interface{}) {
	t := currentThread(true)
	mp := tls.getMap(t)
	if mp != nil {
		mp.set(tls, value)
	} else {
		tls.createMap(t, value)
	}
}

func (tls *inheritableThreadLocalImpl) Remove() {
	t := currentThread(true)
	mp := tls.getMap(t)
	if mp != nil {
		mp.remove(tls)
	}
}

func (tls *inheritableThreadLocalImpl) getMap(t *thread) *threadLocalMap {
	return t.inheritableThreadLocals
}

func (tls *inheritableThreadLocalImpl) createMap(t *thread, firstValue interface{}) {
	mp := &threadLocalMap{}
	mp.set(tls, firstValue)
	t.inheritableThreadLocals = mp
}

func (tls *inheritableThreadLocalImpl) setInitialValue(t *thread) interface{} {
	value := tls.initialValue()
	mp := tls.getMap(t)
	if mp != nil {
		mp.set(tls, value)
	} else {
		tls.createMap(t, value)
	}
	return value
}

func (tls *inheritableThreadLocalImpl) initialValue() interface{} {
	if tls.supplier == nil {
		return nil
	}
	return tls.supplier()
}

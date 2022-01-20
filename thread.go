package routine

import (
	"sync"
	"sync/atomic"
)

var (
	globalMap     atomic.Value // The global thread map (map[int64]*thread)
	globalMapLock sync.Mutex   // The Lock to control accessing of globalMap
)

func init() {
	globalMap.Store(map[int64]*thread{})
}

type thread struct {
	id                      int64
	threadLocals            *threadLocalMap
	inheritableThreadLocals *threadLocalMap
}

func currentThread(create bool) *thread {
	gid := Goid()
	gMap := globalMap.Load().(map[int64]*thread)
	var t *thread
	if t = gMap[gid]; t == nil && create {
		t = &thread{
			id: gid,
		}
		globalMapLock.Lock()
		defer globalMapLock.Unlock()
		oldGMap := globalMap.Load().(map[int64]*thread)
		newGMap := make(map[int64]*thread, len(oldGMap)+1)
		for k, v := range oldGMap {
			newGMap[k] = v
		}
		newGMap[gid] = t
		globalMap.Store(newGMap)
	}
	return t
}

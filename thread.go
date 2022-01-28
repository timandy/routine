package routine

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	segmentSize = 16
)

var (
	globalMap     = [segmentSize]*atomic.Value{} // The global thread map (map[int64]*thread)
	globalMapLock = [segmentSize]*sync.Mutex{}   // The Lock to control accessing of globalMap
)

func init() {
	for idx := 0; idx < segmentSize; idx++ {
		segment := &atomic.Value{}
		segment.Store(map[int64]*thread{})
		globalMap[idx] = segment
		globalMapLock[idx] = &sync.Mutex{}
	}
}

type thread struct {
	id                      int64
	threadLocals            *threadLocalMap
	inheritableThreadLocals *threadLocalMap
}

func hash(goid int64) int {
	idx := int(goid % segmentSize)
	if idx < 0 {
		return -idx
	}
	return idx
}

func currentThread(create bool) *thread {
	goid := Goid()
	idx := hash(goid)
	segment := globalMap[idx]
	segmentMap := segment.Load().(map[int64]*thread)
	var t *thread
	if t = segmentMap[goid]; t == nil && create {
		t = &thread{
			id: goid,
		}
		segmentLock := globalMapLock[idx]
		segmentLock.Lock()
		defer segmentLock.Unlock()
		oldSegMap := segment.Load().(map[int64]*thread)
		newSegMap := make(map[int64]*thread, len(oldSegMap)+1)
		for k, v := range oldSegMap {
			newSegMap[k] = v
		}
		newSegMap[goid] = t
		segment.Store(newSegMap)
		// try restart gc timer if Set for the first time
		if gcTimer[idx] == nil {
			gcTimer[idx] = time.AfterFunc(gCInterval, func() {
				gc(idx)
			})
		}
	}
	return t
}

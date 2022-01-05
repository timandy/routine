package routine

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	globalMap     atomic.Value       // The global threadLocalImpl map (map[int64]*threadLocalMap)
	globalMapLock sync.Mutex         // The Lock to control accessing of globalMap
	gcTimer       *time.Timer        // The timer of globalMap's garbage collector
	gCInterval    = time.Second * 30 // The pre-defined gc interval
)

func init() {
	globalMap.Store(map[int64]*threadLocalMap{})
}

func gcRunning() bool {
	globalMapLock.Lock()
	defer globalMapLock.Unlock()
	return gcTimer != nil
}

type threadLocalMap struct {
	gid    int64
	values []interface{}
}

func (s *threadLocalMap) get(index int) interface{} {
	if index < len(s.values) {
		return s.values[index]
	}
	return nil
}

func (s *threadLocalMap) set(index int, value interface{}) interface{} {
	if index < len(s.values) {
		oldValue := s.values[index]
		s.values[index] = value
		return oldValue
	}

	newCapacity := index
	newCapacity |= newCapacity >> 1
	newCapacity |= newCapacity >> 2
	newCapacity |= newCapacity >> 4
	newCapacity |= newCapacity >> 8
	newCapacity |= newCapacity >> 16
	newCapacity++

	newValues := make([]interface{}, newCapacity)
	copy(newValues, s.values)
	newValues[index] = value
	s.values = newValues
	return nil
}

func (s *threadLocalMap) remove(index int) interface{} {
	if index < len(s.values) {
		oldValue := s.values[index]
		s.values[index] = nil
		return oldValue
	}
	return nil
}

func (s *threadLocalMap) clear() {
	s.values = []interface{}{}
}

type threadLocalImpl struct {
	id int
}

func (t *threadLocalImpl) Get() interface{} {
	s := getMap(false)
	if s == nil {
		return nil
	}
	return s.get(t.id)
}

func (t *threadLocalImpl) Set(value interface{}) interface{} {
	s := getMap(true)
	oldValue := s.set(t.id, value)

	// try restart gc timer if Set for the first time
	if oldValue == nil {
		globalMapLock.Lock()
		if gcTimer == nil {
			gcTimer = time.AfterFunc(gCInterval, gc)
		}
		globalMapLock.Unlock()
	}
	return oldValue
}

func (t *threadLocalImpl) Remove() interface{} {
	s := getMap(false)
	if s == nil {
		return nil
	}
	return s.remove(t.id)
}

// getMap load the threadLocalMap of current goroutine.
func getMap(create bool) *threadLocalMap {
	gid := Goid()
	storeMap := globalMap.Load().(map[int64]*threadLocalMap)
	var s *threadLocalMap
	if s = storeMap[gid]; s == nil && create {
		globalMapLock.Lock()
		oldStoreMap := globalMap.Load().(map[int64]*threadLocalMap)
		if s = oldStoreMap[gid]; s == nil {
			s = &threadLocalMap{
				gid:    gid,
				values: make([]interface{}, 8),
			}
			newStoreMap := make(map[int64]*threadLocalMap, len(oldStoreMap)+1)
			for k, v := range oldStoreMap {
				newStoreMap[k] = v
			}
			newStoreMap[gid] = s
			globalMap.Store(newStoreMap)
		}
		globalMapLock.Unlock()
	}
	return s
}

// gc clear all data of dead goroutine.
func gc() {
	globalMapLock.Lock()
	defer globalMapLock.Unlock()

	// load all alive goids
	gids := AllGoids()
	gidMap := make(map[int64]struct{}, len(gids))
	for _, gid := range gids {
		gidMap[gid] = struct{}{}
	}

	// scan global storeMap check the dead and live threadLocalMap count.
	var storeMap = globalMap.Load().(map[int64]*threadLocalMap)
	var liveCnt int
	for gid := range storeMap {
		if _, ok := gidMap[gid]; ok {
			liveCnt++
		}
	}

	// clean dead threadLocalMap of dead goroutine if needed.
	if liveCnt != len(storeMap) {
		newStoreMap := make(map[int64]*threadLocalMap, liveCnt)
		for gid, s := range storeMap {
			if _, ok := gidMap[gid]; ok {
				newStoreMap[gid] = s
			}
		}
		globalMap.Store(newStoreMap)
	}

	// setup next round timer if need. TODO it's ok?
	if liveCnt > 0 {
		gcTimer.Reset(gCInterval)
	} else {
		gcTimer = nil
	}
}

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

func (mp *threadLocalMap) get(index int) interface{} {
	if index < len(mp.values) {
		return mp.values[index]
	}
	return nil
}

func (mp *threadLocalMap) set(index int, value interface{}) interface{} {
	if index < len(mp.values) {
		oldValue := mp.values[index]
		mp.values[index] = value
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
	copy(newValues, mp.values)
	newValues[index] = value
	mp.values = newValues
	return nil
}

func (mp *threadLocalMap) remove(index int) interface{} {
	if index < len(mp.values) {
		oldValue := mp.values[index]
		mp.values[index] = nil
		return oldValue
	}
	return nil
}

func (mp *threadLocalMap) clear() {
	mp.values = []interface{}{}
}

type threadLocalImpl struct {
	id int
}

func (tls *threadLocalImpl) Get() interface{} {
	mp := getMap(false)
	if mp == nil {
		return nil
	}
	return mp.get(tls.id)
}

func (tls *threadLocalImpl) Set(value interface{}) interface{} {
	mp := getMap(true)
	oldValue := mp.set(tls.id, value)

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

func (tls *threadLocalImpl) Remove() interface{} {
	mp := getMap(false)
	if mp == nil {
		return nil
	}
	return mp.remove(tls.id)
}

// getMap load the threadLocalMap of current goroutine.
func getMap(create bool) *threadLocalMap {
	gid := Goid()
	gMap := globalMap.Load().(map[int64]*threadLocalMap)
	var lMap *threadLocalMap
	if lMap = gMap[gid]; lMap == nil && create {
		globalMapLock.Lock()
		oldGMap := globalMap.Load().(map[int64]*threadLocalMap)
		if lMap = oldGMap[gid]; lMap == nil {
			lMap = &threadLocalMap{
				gid:    gid,
				values: make([]interface{}, 8),
			}
			newGMap := make(map[int64]*threadLocalMap, len(oldGMap)+1)
			for k, v := range oldGMap {
				newGMap[k] = v
			}
			newGMap[gid] = lMap
			globalMap.Store(newGMap)
		}
		globalMapLock.Unlock()
	}
	return lMap
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

	// scan globalMap check the dead and live threadLocalMap count.
	var gMap = globalMap.Load().(map[int64]*threadLocalMap)
	var liveCnt int
	for gid := range gMap {
		if _, ok := gidMap[gid]; ok {
			liveCnt++
		}
	}

	// clean dead threadLocalMap of dead goroutine if needed.
	if liveCnt != len(gMap) {
		newGMap := make(map[int64]*threadLocalMap, liveCnt)
		for gid, lMap := range gMap {
			if _, ok := gidMap[gid]; ok {
				newGMap[gid] = lMap
			}
		}
		globalMap.Store(newGMap)
	}

	// setup next round timer if need. TODO it's ok?
	if liveCnt > 0 {
		gcTimer.Reset(gCInterval)
	} else {
		gcTimer = nil
	}
}

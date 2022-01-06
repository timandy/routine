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

type entry struct {
	value interface{}
}

type threadLocalMap struct {
	gid     int64
	entries []*entry
}

func (mp *threadLocalMap) getEntry(key *threadLocalImpl) *entry {
	index := key.id
	if index < len(mp.entries) {
		return mp.entries[index]
	}
	return nil
}

func (mp *threadLocalMap) set(key *threadLocalImpl, value interface{}) bool {
	index := key.id
	if index < len(mp.entries) {
		e := mp.entries[index]
		if e == nil {
			mp.entries[index] = &entry{value: value}
			return true
		}
		e.value = value
		return false
	}

	newCapacity := index
	newCapacity |= newCapacity >> 1
	newCapacity |= newCapacity >> 2
	newCapacity |= newCapacity >> 4
	newCapacity |= newCapacity >> 8
	newCapacity |= newCapacity >> 16
	newCapacity++

	newEntries := make([]*entry, newCapacity)
	copy(newEntries, mp.entries)
	newEntries[index] = &entry{value: value}
	mp.entries = newEntries
	return true
}

func (mp *threadLocalMap) remove(key *threadLocalImpl) {
	index := key.id
	if index < len(mp.entries) {
		mp.entries[index] = nil
	}
}

func (mp *threadLocalMap) clear() {
	mp.entries = []*entry{}
}

type threadLocalImpl struct {
	id int
}

func (tls *threadLocalImpl) Get() interface{} {
	mp := getMap(false)
	if mp == nil {
		return nil
	}
	e := mp.getEntry(tls)
	if e == nil {
		return nil
	}
	return e.value
}

func (tls *threadLocalImpl) Set(value interface{}) {
	mp := getMap(true)
	notExists := mp.set(tls, value)

	// try restart gc timer if Set for the first time
	if notExists {
		globalMapLock.Lock()
		if gcTimer == nil {
			gcTimer = time.AfterFunc(gCInterval, gc)
		}
		globalMapLock.Unlock()
	}
}

func (tls *threadLocalImpl) Remove() {
	mp := getMap(false)
	if mp == nil {
		return
	}
	mp.remove(tls)
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
				gid:     gid,
				entries: make([]*entry, 8),
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

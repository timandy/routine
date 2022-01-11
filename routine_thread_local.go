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

func (mp *threadLocalMap) set(key *threadLocalImpl, value interface{}) {
	index := key.id
	if index < len(mp.entries) {
		e := mp.entries[index]
		if e == nil {
			mp.entries[index] = &entry{value: value}
			// try restart gc timer if Set for the first time
			gcTimerStart()
		} else {
			e.value = value
		}
		return
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
	// try restart gc timer if Set for the first time
	gcTimerStart()
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
	id       int
	supplier func() interface{}
}

func (tls *threadLocalImpl) Get() interface{} {
	mp := getMap(true)
	e := mp.getEntry(tls)
	if e != nil {
		return e.value
	}
	return tls.setInitialValue(mp)
}

func (tls *threadLocalImpl) Set(value interface{}) {
	mp := getMap(true)
	mp.set(tls, value)
}

func (tls *threadLocalImpl) Remove() {
	mp := getMap(false)
	if mp == nil {
		return
	}
	mp.remove(tls)
}

func (tls *threadLocalImpl) setInitialValue(mp *threadLocalMap) interface{} {
	value := tls.initialValue()
	mp.set(tls, value)
	return value
}

func (tls *threadLocalImpl) initialValue() interface{} {
	if tls.supplier == nil {
		return nil
	}
	return tls.supplier()
}

// getMap load the threadLocalMap of current goroutine.
func getMap(create bool) *threadLocalMap {
	gid := Goid()
	gMap := globalMap.Load().(map[int64]*threadLocalMap)
	var lMap *threadLocalMap
	if lMap = gMap[gid]; lMap == nil && create {
		lMap = &threadLocalMap{
			gid:     gid,
			entries: make([]*entry, 8),
		}
		globalMapLock.Lock()
		defer globalMapLock.Unlock()
		oldGMap := globalMap.Load().(map[int64]*threadLocalMap)
		newGMap := make(map[int64]*threadLocalMap, len(oldGMap)+1)
		for k, v := range oldGMap {
			newGMap[k] = v
		}
		newGMap[gid] = lMap
		globalMap.Store(newGMap)
	}
	return lMap
}

// gcRunning if gcTimer is not nil return true, else return false
func gcRunning() bool {
	globalMapLock.Lock()
	defer globalMapLock.Unlock()
	return gcTimer != nil
}

// gcTimerStart make sure gcTimer is not nil
func gcTimerStart() {
	globalMapLock.Lock()
	defer globalMapLock.Unlock()
	if gcTimer == nil {
		gcTimer = time.AfterFunc(gCInterval, gc)
	}
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

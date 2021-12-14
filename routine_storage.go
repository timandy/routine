package routine

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	storages          atomic.Value       // The global storage map (map[int64]*store)
	storageLock       sync.Mutex         // The Lock to control accessing of storages
	storageGCTimer    *time.Timer        // The timer of storage's garbage collector
	storageGCInterval = time.Second * 30 // The pre-defined gc interval
)

func init() {
	storages.Store(map[int64]*store{})
}

func gcRunning() bool {
	storageLock.Lock()
	defer storageLock.Unlock()
	return storageGCTimer != nil
}

type store struct {
	gid    int64
	values []interface{}
}

func (s *store) get(index int) interface{} {
	if index < len(s.values) {
		return s.values[index]
	}
	return nil
}

func (s *store) set(index int, value interface{}) interface{} {
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

func (s *store) remove(index int) interface{} {
	if index < len(s.values) {
		oldValue := s.values[index]
		s.values[index] = nil
		return oldValue
	}
	return nil
}

func (s *store) clear() {
	s.values = []interface{}{}
}

type storage struct {
	id int
}

func (t *storage) Get() (value interface{}) {
	s := loadCurrentStore(false)
	if s == nil {
		return nil
	}
	value = s.get(t.id)
	return
}

func (t *storage) Set(value interface{}) (oldValue interface{}) {
	s := loadCurrentStore(true)
	oldValue = s.set(t.id, value)

	// try restart gc timer if Set for the first time
	if oldValue == nil {
		storageLock.Lock()
		if storageGCTimer == nil {
			storageGCTimer = time.AfterFunc(storageGCInterval, clearDeadStore)
		}
		storageLock.Unlock()
	}
	return
}

func (t *storage) Remove() (oldValue interface{}) {
	s := loadCurrentStore(false)
	if s == nil {
		return nil
	}
	oldValue = s.remove(t.id)
	return
}

// loadCurrentStore load the store of current goroutine.
func loadCurrentStore(create bool) (s *store) {
	gid := Goid()
	storeMap := storages.Load().(map[int64]*store)
	if s = storeMap[gid]; s == nil && create {
		storageLock.Lock()
		oldStoreMap := storages.Load().(map[int64]*store)
		if s = oldStoreMap[gid]; s == nil {
			s = &store{
				gid:    gid,
				values: make([]interface{}, 8),
			}
			newStoreMap := make(map[int64]*store, len(oldStoreMap)+1)
			for k, v := range oldStoreMap {
				newStoreMap[k] = v
			}
			newStoreMap[gid] = s
			storages.Store(newStoreMap)
		}
		storageLock.Unlock()
	}
	return
}

// clearDeadStore clear all data of dead goroutine.
func clearDeadStore() {
	storageLock.Lock()
	defer storageLock.Unlock()

	// load all alive goids
	gids := AllGoids()
	gidMap := make(map[int64]struct{}, len(gids))
	for _, gid := range gids {
		gidMap[gid] = struct{}{}
	}

	// scan global storeMap check the dead and live store count.
	var storeMap = storages.Load().(map[int64]*store)
	var deadCnt, liveCnt int
	for gid := range storeMap {
		if _, ok := gidMap[gid]; ok {
			liveCnt++
		} else {
			deadCnt++
		}
	}

	// clean dead store of dead goroutine if need.
	if deadCnt > 0 {
		newStoreMap := make(map[int64]*store, len(storeMap)-deadCnt)
		for gid, s := range storeMap {
			if _, ok := gidMap[gid]; ok {
				newStoreMap[gid] = s
			}
		}
		storages.Store(newStoreMap)
	}

	// setup next round timer if need. TODO it's ok?
	if liveCnt > 0 {
		storageGCTimer.Reset(storageGCInterval)
	} else {
		storageGCTimer = nil
	}
}

package routine

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
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
	count  uint32
	values map[uintptr]interface{}
}

type storage struct {
}

func (t *storage) Get() (v interface{}) {
	s := loadCurrentStore()
	id := uintptr(unsafe.Pointer(t))
	return s.values[id]
}

func (t *storage) Set(v interface{}) (oldValue interface{}) {
	s := loadCurrentStore()
	id := uintptr(unsafe.Pointer(t))
	oldValue = s.values[id]
	s.values[id] = v
	atomic.StoreUint32(&s.count, uint32(len(s.values)))

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

func (t *storage) Del() (v interface{}) {
	s := loadCurrentStore()
	id := uintptr(unsafe.Pointer(t))
	v = s.values[id]
	delete(s.values, id)
	atomic.StoreUint32(&s.count, uint32(len(s.values)))
	return
}

func (t *storage) Clear() {
	s := loadCurrentStore()
	s.values = map[uintptr]interface{}{}
	atomic.StoreUint32(&s.count, 0)
}

// loadCurrentStore load the store of current goroutine.
func loadCurrentStore() (s *store) {
	gid := Goid()
	storeMap := storages.Load().(map[int64]*store)
	if s = storeMap[gid]; s == nil {
		storageLock.Lock()
		oldStoreMap := storages.Load().(map[int64]*store)
		if s = oldStoreMap[gid]; s == nil {
			s = &store{
				gid:    gid,
				values: map[uintptr]interface{}{},
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
	for id, s := range storeMap {
		if _, ok := gidMap[id]; ok {
			if atomic.LoadUint32(&s.count) > 0 {
				liveCnt++
			}
			liveCnt++
		} else {
			deadCnt++
		}
	}

	// clean dead store of dead goroutine if need.
	if deadCnt > 0 {
		newStoreMap := make(map[int64]*store, len(storeMap)-deadCnt)
		for id, s := range storeMap {
			if _, ok := gidMap[id]; ok {
				newStoreMap[id] = s
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

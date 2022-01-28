package routine

import "time"

var (
	gcTimer    *time.Timer        // The timer of globalMap's garbage collector
	gCInterval = time.Second * 30 // The pre-defined gc interval
)

// gcRunning if gcTimer is not nil return true, else return false
func gcRunning() bool {
	globalMapLock.Lock()
	defer globalMapLock.Unlock()
	return gcTimer != nil
}

// gc clear all data of dead goroutine.
func gc() {
	globalMapLock.Lock()
	defer globalMapLock.Unlock()

	// load all alive goids
	goids := AllGoids()
	goidMap := make(map[int64]struct{}, len(goids))
	for _, goid := range goids {
		goidMap[goid] = struct{}{}
	}

	// compute how many thread instances are there *at most* after GC.
	gMap := globalMap.Load().(map[int64]*thread)
	gMapLen := len(gMap)
	liveCnt := len(goidMap)
	if liveCnt > gMapLen {
		liveCnt = gMapLen
	}

	// clean dead thread of dead goroutine.
	newGMap := make(map[int64]*thread, liveCnt)
	for goid, t := range gMap {
		if _, ok := goidMap[goid]; ok {
			newGMap[goid] = t
		}
	}
	globalMap.Store(newGMap)

	// setup next round timer if needed.
	if len(newGMap) > 0 {
		gcTimer.Reset(gCInterval)
	} else {
		gcTimer = nil
	}
}

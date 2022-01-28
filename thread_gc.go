package routine

import "time"

var (
	gcTimer    = [segmentSize]*time.Timer{} // The timers of globalMap's garbage collector
	gCInterval = time.Second * 30           // The pre-defined gc interval
)

// gc clear all data of dead goroutine.
func gc(idx int) {
	segmentLock := globalMapLock[idx]
	segmentLock.Lock()
	defer segmentLock.Unlock()

	// load all alive goids
	goids := AllGoids()
	goidMap := make(map[int64]struct{}, len(goids))
	for _, goid := range goids {
		if hash(goid) == idx {
			goidMap[goid] = struct{}{}
		}
	}
	// compute how many thread instances are there *at most* after GC.
	segment := globalMap[idx]
	segmentMap := segment.Load().(map[int64]*thread)
	segmentMapLen := len(segmentMap)
	liveCnt := len(goidMap)
	if liveCnt > segmentMapLen {
		liveCnt = segmentMapLen
	}

	// clean dead thread of dead goroutine.
	newSegMap := make(map[int64]*thread, liveCnt)
	for goid, t := range segmentMap {
		if _, ok := goidMap[goid]; ok {
			newSegMap[goid] = t
		}
	}
	segment.Store(newSegMap)

	// setup next round timer if needed.
	if len(newSegMap) > 0 {
		gcTimer[idx].Reset(gCInterval)
	} else {
		gcTimer[idx] = nil
	}
}

package routine

import (
	"runtime"
	"strings"
	"unsafe"
)

var (
	goidOffset    uintptr
	goidOffsetDic = map[string]uintptr{
		"go1.13": 152,
		"go1.14": 152,
		"go1.15": 152,
		"go1.16": 152,
		"go1.17": 152,
		"go1.18": 152,
	}
)

func init() {
	var offset uintptr
	version := runtime.Version()
	for k, v := range goidOffsetDic {
		if k == version || strings.HasPrefix(version, k) {
			offset = v
			break
		}
	}
	goidOffset = offset
}

// getGoidByNative parse the current goroutine's id from g.
// This function could be very fast(like 1ns/op), but it may be failed.
func getGoidByNative() (int64, bool) {
	if goidOffset == 0 {
		return 0, false
	}
	gp := getg()
	if gp == nil {
		return 0, false
	}
	goid := (*int64)(add(gp, goidOffset))
	if goid == nil {
		return 0, false
	}
	return *goid, true
}

// getGoidByStack parse the current goroutine's id from caller stack.
// This function could be very slow(like 3000us/op), but it's very safe.
func getGoidByStack() int64 {
	return int64(curGoroutineID())
}

//add pointer addition operation.
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

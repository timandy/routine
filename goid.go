package routine

import (
	"fmt"
	"github.com/timandy/routine/g"
	"runtime"
	"strings"
	"unsafe"
)

var (
	goidOffset    uintptr
	anchor        = []byte("goroutine ")
	goidOffsetDic = map[string]int64{
		"go1.13": 152,
		"go1.14": 152,
		"go1.15": 152,
		"go1.16": 152,
		"go1.17": 152,
	}
)

func init() {
	var off int64
	version := runtime.Version()
	for k, v := range goidOffsetDic {
		if version == k || strings.HasPrefix(version, k) {
			off = v
			break
		}
	}
	goidOffset = uintptr(off)
}

// getGoidByNative parse the current goroutine's id from G.
// This function could be very fast(like 1ns/op), but it may be failed.
func getGoidByNative() (int64, bool) {
	gp := g.G()
	goid := findGoidPointer(gp)
	if goid == nil {
		return 0, false
	}
	return *goid, true
}

// getGoidByStack parse the current goroutine's id from caller stack.
// This function could be very slow(like 3000us/op), but it's very safe.
func getGoidByStack() int64 {
	buf := traceTiny()
	goid, _ := findNextGoid(buf, 0)
	return goid
}

// getAllGoidByStack retrieve all goid through stack;
// This function could be very inefficient, but it's very safe.
func getAllGoidByStack() []int64 {
	buf := traceAllStack()
	// parse all goids
	goids := make([]int64, 0, 100)
	for i := 0; i < len(buf); {
		goid, off := findNextGoid(buf, i)
		if goid > 0 {
			goids = append(goids, goid)
		}
		i = off
	}
	return goids
}

// foreachGoidByStack run a func for each goroutine's goid through stack;
// This function could be very inefficient, but it's very safe.
func foreachGoidByStack(fun func(goid int64)) {
	buf := traceAllStack()
	for i := 0; i < len(buf); {
		goid, off := findNextGoid(buf, i)
		if goid > 0 {
			fun(goid)
		}
		i = off
	}
}

// Return the pointer of its goid through the pointer of the g structure.
func findGoidPointer(gp unsafe.Pointer) *int64 {
	if goidOffset == 0 || gp == nil {
		return nil
	}
	return (*int64)(unsafe.Pointer(uintptr(gp) + goidOffset))
}

// Find the next goid from buf[off:]
func findNextGoid(buf []byte, off int) (goid int64, next int) {
	i := off
	hit := false
	// skip to anchor
	anc := anchor
	bufLen := len(buf)
	ancLen := len(anc)
	for stop := bufLen - ancLen; i < stop; {
		if buf[i] == anc[0] && buf[i+1] == anc[1] && buf[i+2] == anc[2] && buf[i+3] == anc[3] &&
			buf[i+4] == anc[4] && buf[i+5] == anc[5] && buf[i+6] == anc[6] &&
			buf[i+7] == anc[7] && buf[i+8] == anc[8] && buf[i+9] == anc[9] {
			hit = true
			i += ancLen
			break
		}
		for ; i < bufLen && buf[i] != '\n'; i++ {
		}
		i++
	}
	// return if not hit
	if !hit {
		return 0, bufLen
	}
	// extract goid
	done := false
	for ; i < bufLen && !done; i++ {
		switch buf[i] {
		case '0':
			goid *= 10
		case '1':
			goid = goid*10 + 1
		case '2':
			goid = goid*10 + 2
		case '3':
			goid = goid*10 + 3
		case '4':
			goid = goid*10 + 4
		case '5':
			goid = goid*10 + 5
		case '6':
			goid = goid*10 + 6
		case '7':
			goid = goid*10 + 7
		case '8':
			goid = goid*10 + 8
		case '9':
			goid = goid*10 + 9
		case ' ':
			done = true
			break
		default:
			goid = 0
			fmt.Println("Should never be here, any bug happens")
		}
	}
	next = i
	return
}

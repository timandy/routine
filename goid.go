package routine

import (
	"fmt"
	"github.com/timandy/routine/g"
	"runtime"
	"strings"
	"unsafe"
)

const (
	ptrSize = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const
)

var (
	goidOffset     uintptr
	anchor         = []byte("goroutine ")
	goidOffsetDict = map[string]int64{
		"go1.12": 152,
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
	for k, v := range goidOffsetDict {
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
	if goidOffset == 0 {
		return 0, false
	}
	tmp := g.G()
	if tmp == nil {
		return 0, false
	}
	p := (*int64)(unsafe.Pointer(uintptr(tmp) + goidOffset))
	if p == nil {
		return 0, false
	}
	return *p, true
}

// getGoidByStack parse the current goroutine's id from caller stack.
// This function could be very slow(like 3000us/op), but it's very safe.
func getGoidByStack() int64 {
	buf := traceTiny()
	goid, _ := findNextGoid(buf, 0)
	return goid
}

// getAllGoidByStack find all goid through stack;
// WARNING: This function could be very inefficient; This method is not thread safe
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

// Find the next goid from `buf[off:]`
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
	var done bool
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

package routine

import (
	"fmt"
	"github.com/timandy/routine/g"
	"runtime"
	"strings"
	"sync"
	"unsafe"
)

const (
	ptrSize   = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const
	stackSize = 1024
)

var (
	goidOffset      uintptr
	allStackBuf     []byte
	allStackBufSize = stackSize * 1024
	anchor          = []byte("goroutine ")
	stackBufPool    = sync.Pool{
		New: func() interface{} {
			buf := make([]byte, 64)
			return &buf
		},
	}
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

// getGoidByStack parse the current goroutine's id from caller stack.
// This function could be very slow(like 3000us/op), but it's very safe.
func getGoidByStack() (goid int64) {
	bp := stackBufPool.Get().(*[]byte)
	defer stackBufPool.Put(bp)

	b := *bp
	b = b[:runtime.Stack(b, false)]
	goid, _ = findNextGoid(b, 0)
	return
}

// getGoidByNative parse the current goroutine's id from G.
// This function could be very fast(like 1ns/op), but it could be failed.
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

// getAllGoidByStack find all goid through stack; WARNING: This function could be very inefficient; This method is not thread safe
func getAllGoidByStack() (goids []int64) {
	buf := readAllStackBuf()
	defer releaseAllStackBuf()
	// parse all goids
	goids = make([]int64, 0, 100)
	for i := 0; i < len(buf); {
		goid, off := findNextGoid(buf, i)
		if goid > 0 {
			goids = append(goids, goid)
		}
		i = off
	}
	return
}

func readStackBuf() []byte {
	stackBuf := make([]byte, stackSize)
	written := runtime.Stack(stackBuf, false)
	for written >= len(stackBuf) {
		stackBuf = make([]byte, len(stackBuf)<<1)
		written = runtime.Stack(stackBuf, false)
	}
	return stackBuf[0:written]
}

// Read all stack info into a buf
func readAllStackBuf() []byte {
	if allStackBuf == nil {
		allStackBuf = make([]byte, allStackBufSize)
	}
	written := runtime.Stack(allStackBuf, true)
	for written >= len(allStackBuf) {
		allStackBuf = make([]byte, len(allStackBuf)<<1)
		written = runtime.Stack(allStackBuf, true)
	}
	return allStackBuf[0:written]
}

// Release stack buf when it is too large
func releaseAllStackBuf() {
	if allStackBuf == nil || len(allStackBuf) == allStackBufSize {
		return
	}
	allStackBuf = nil
}

// Find the next goid from `buf[off:]`
func findNextGoid(buf []byte, off int) (goid int64, next int) {
	i := off
	hit := false
	// skip to anchor
	acr := anchor
	for sb := len(buf) - len(acr); i < sb; {
		if buf[i] == acr[0] && buf[i+1] == acr[1] && buf[i+2] == acr[2] && buf[i+3] == acr[3] &&
			buf[i+4] == acr[4] && buf[i+5] == acr[5] && buf[i+6] == acr[6] &&
			buf[i+7] == acr[7] && buf[i+8] == acr[8] && buf[i+9] == acr[9] {
			hit = true
			i += len(acr)
			break
		}
		for ; i < len(buf) && buf[i] != '\n'; i++ {
		}
		i++
	}
	// return if not hit
	if !hit {
		return 0, len(buf)
	}
	// extract goid
	var done bool
	for ; i < len(buf) && !done; i++ {
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
			fmt.Println("should never be here, any bug happens")
		}
	}
	next = i
	return
}

package routine

import (
	"runtime"
)

const (
	tinyStackSize = 64
	curStackSize  = 1024
	allStackSize  = 1024 * 1024
)

func traceTiny() []byte {
	buf := make([]byte, tinyStackSize)
	written := runtime.Stack(buf, false)
	return buf[0:written]
}

func traceStack() []byte {
	buf := make([]byte, curStackSize)
	written := runtime.Stack(buf, false)
	for written >= len(buf) {
		buf = make([]byte, len(buf)<<1)
		written = runtime.Stack(buf, false)
	}
	return buf[0 : written-1] //remove last \n
}

func traceAllStack() []byte {
	buf := make([]byte, allStackSize)
	written := runtime.Stack(buf, true)
	for written >= len(buf) {
		buf = make([]byte, len(buf)<<1)
		written = runtime.Stack(buf, true)
	}
	return buf[0 : written-1] //remove last \n
}

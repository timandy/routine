package routine

import "runtime"

const stackSize = 1024

func traceStack() []byte {
	buf := make([]byte, stackSize)
	written := runtime.Stack(buf, false)
	for written >= len(buf) {
		buf = make([]byte, len(buf)<<1)
		written = runtime.Stack(buf, false)
	}
	return buf[0 : written-1] //remove last \n
}

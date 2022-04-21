package routine

import "runtime"

func traceStack() []byte {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	for n >= len(buf) {
		buf = make([]byte, len(buf)<<1)
		n = runtime.Stack(buf, false)
	}
	return buf[:n-1] //remove last \n
}

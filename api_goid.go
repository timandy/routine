package routine

import "fmt"

// Goid return the current goroutine's unique id.
// It will try to get goid by native cgo/asm for better performance,
// and could parse goid from stack for fail over supporting.
func Goid() int64 {
	if goid, success := getGoidByNative(); success {
		return goid
	}
	fmt.Println("[WARNING] Unable to get goid from runtime natively, now fall over to stack info which will be very inefficient!!!")
	return getGoidByStack()
}

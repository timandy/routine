package routine

import "fmt"

// Goid return the current goroutine's unique id.
// It will try get gid by native cgo/asm for better performance,
// and could parse gid from stack for failover supporting.
func Goid() int64 {
	if goid, success := getGoidByNative(); success {
		return goid
	}
	return getGoidByStack()
}

// AllGoids return all goroutine's goid in the current golang process.
// It will try load all goid from runtime natively for better performance,
// and fallover to runtime.Stack, which is realy inefficient.
func AllGoids() []int64 {
	if goids, err := getAllGoidByNative(); err == nil {
		return goids
	}
	fmt.Println("[WARNING] cannot get all goids from runtime natively, now fall over to stack info, this will be very inefficient!!!")
	return getAllGoidByStack()
}

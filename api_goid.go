package routine

import "fmt"

// Goid return the current goroutine's unique id.
// It will try to get goid by native cgo/asm for better performance,
// and could parse goid from stack for fail over supporting.
func Goid() int64 {
	if goid, success := getGoidByNative(); success {
		return goid
	}
	fmt.Println("[WARNING] cannot get goid from runtime natively, now fall over to stack info, this will be very inefficient!!!")
	return getGoidByStack()
}

// AllGoids return all goroutine's goid in the current golang process.
// It will try load all goid from runtime natively for better performance,
// and fall over to runtime.Stack, which is very inefficient.
// Addition of new Gs during execution, which may be missed.
func AllGoids() []int64 {
	if goids, success := getAllGoidByNative(); success {
		return goids
	}
	fmt.Println("[WARNING] cannot get all goids from runtime natively, now fall over to stack info, this will be very inefficient!!!")
	return getAllGoidByStack()
}

// ForeachGoid run a func for each goroutine's goid in the current golang process.
// It will try load all goid from runtime natively for better performance,
// and fall over to runtime.Stack, which is very inefficient.
// Addition of new Gs during execution, which may be missed.
func ForeachGoid(fun func(goid int64)) {
	if success := foreachGoidByNative(fun); success {
		return
	}
	fmt.Println("[WARNING] cannot run func for each goid from runtime natively, now fall over to stack info, this will be very inefficient!!!")
	foreachGoidByStack(fun)
}

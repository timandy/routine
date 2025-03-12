//go:build routinex

package routine

import "unsafe"

type thread struct {
	threadLocals            *threadLocalMap
	inheritableThreadLocals *threadLocalMap
}

// currentThread returns a pointer to the currently executing goroutine's thread struct.
//
//go:norace
//go:nocheckptr
func currentThread(create bool) *thread {
	gp := getg()
	return (*thread)(add(unsafe.Pointer(gp), offsetThreadLocals))
}

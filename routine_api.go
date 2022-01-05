package routine

import (
	"fmt"
	"sync/atomic"
)

// ThreadLocal provides goroutine-local variables.
type ThreadLocal interface {
	// Get returns the value in the current goroutine's local threadLocalImpl, if it was set before.
	Get() interface{}

	// Set copy the value into the current goroutine's local threadLocalImpl, and return the old value.
	Set(value interface{}) interface{}

	// Remove delete the value from the current goroutine's local threadLocalImpl, and return it.
	Remove() interface{}
}

// Clear clean up all context variables of the current coroutine.
func Clear() {
	s := getMap(false)
	if s == nil {
		return
	}
	s.clear()
}

// ImmutableContext represents all local values of one goroutine.
type ImmutableContext struct {
	gid    int64
	values []interface{}
}

// Go starts a new goroutine, and copy all local values from current goroutine.
func Go(f func()) {
	ic := BackupContext()
	go func() {
		RestoreContext(ic)
		f()
	}()
}

// BackupContext copy all local values into an ImmutableContext instance.
func BackupContext() *ImmutableContext {
	s := getMap(false)
	if s == nil || s.values == nil {
		return nil
	}
	data := make([]interface{}, len(s.values))
	copy(data, s.values)
	return &ImmutableContext{gid: s.gid, values: data}
}

// RestoreContext load the specified ImmutableContext instance into the local threadLocalImpl of current goroutine.
func RestoreContext(ic *ImmutableContext) {
	if ic == nil || ic.values == nil {
		Clear()
		return
	}
	icLength := len(ic.values)
	s := getMap(true)
	if len(s.values) != icLength {
		s.values = make([]interface{}, icLength)
	}
	copy(s.values, ic.values)
}

var threadLocalIndex int32 = -1

// NewThreadLocal create and return a new ThreadLocal instance.
func NewThreadLocal() ThreadLocal {
	return &threadLocalImpl{id: int(atomic.AddInt32(&threadLocalIndex, 1))}
}

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

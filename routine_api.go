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
	Set(value interface{})

	// Remove delete the value from the current goroutine's local threadLocalImpl, and return it.
	Remove()
}

// Clear clean up all context variables of the current coroutine.
func Clear() {
	mp := getMap(false)
	if mp == nil {
		return
	}
	mp.clear()
}

// ImmutableContext represents all local entries of one goroutine.
type ImmutableContext struct {
	entries []*entry
}

// Go starts a new goroutine, and copy all local entries from current goroutine.
func Go(f func()) {
	ic := BackupContext()
	go func() {
		RestoreContext(ic)
		f()
	}()
}

// BackupContext copy all local entries into an ImmutableContext instance.
func BackupContext() *ImmutableContext {
	mp := getMap(false)
	if mp == nil || mp.entries == nil {
		return nil
	}
	entries := make([]*entry, len(mp.entries))
	copy(entries, mp.entries)
	return &ImmutableContext{entries: entries}
}

// RestoreContext load the specified ImmutableContext instance into the local threadLocalImpl of current goroutine.
func RestoreContext(ic *ImmutableContext) {
	if ic == nil || ic.entries == nil {
		Clear()
		return
	}
	icLength := len(ic.entries)
	mp := getMap(true)
	if len(mp.entries) != icLength {
		mp.entries = make([]*entry, icLength)
	}
	copy(mp.entries, ic.entries)
}

var threadLocalIndex int32 = -1

// NewThreadLocal create and return a new ThreadLocal instance.
func NewThreadLocal() ThreadLocal {
	return &threadLocalImpl{id: int(atomic.AddInt32(&threadLocalIndex, 1))}
}

// NewThreadLocalWithInitial create and return a new ThreadLocal instance. The initial value is determined by invoking the supplier method.
func NewThreadLocalWithInitial(supplier func() interface{}) ThreadLocal {
	return &threadLocalImpl{id: int(atomic.AddInt32(&threadLocalIndex, 1)), supplier: supplier}
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

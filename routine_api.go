package routine

import (
	"fmt"
	"sync/atomic"
)

// LocalStorage provides goroutine-local variables.
type LocalStorage interface {
	// Get returns the value in the current goroutine's local storage, if it was set before.
	Get() interface{}

	// Set copy the value into the current goroutine's local storage, and return the old value.
	Set(value interface{}) interface{}

	// Remove delete the value from the current goroutine's local storage, and return it.
	Remove() interface{}
}

// Clear clean up all context variables of the current coroutine.
func Clear() {
	s := loadCurrentStore(false)
	if s == nil {
		return
	}
	s.clear()
}

// ImmutableContext represents all local storages of one goroutine.
type ImmutableContext struct {
	gid    int64
	values []interface{}
}

// Go start an new goroutine, and copy all local storages from current goroutine.
func Go(f func()) {
	ic := BackupContext()
	go func() {
		RestoreContext(ic)
		f()
	}()
}

// BackupContext copy all local storages into an ImmutableContext instance.
func BackupContext() *ImmutableContext {
	s := loadCurrentStore(false)
	if s == nil || s.values == nil {
		return nil
	}
	data := make([]interface{}, len(s.values))
	copy(data, s.values)
	return &ImmutableContext{gid: s.gid, values: data}
}

// RestoreContext load the specified ImmutableContext instance into the local storage of current goroutine.
func RestoreContext(ic *ImmutableContext) {
	if ic == nil || ic.values == nil {
		Clear()
		return
	}
	icLength := len(ic.values)
	s := loadCurrentStore(true)
	if len(s.values) != icLength {
		s.values = make([]interface{}, icLength)
	}
	copy(s.values, ic.values)
}

var storageIndex int32 = -1

// NewLocalStorage create and return a new LocalStorage instance.
func NewLocalStorage() LocalStorage {
	return &storage{id: int(atomic.AddInt32(&storageIndex, 1))}
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

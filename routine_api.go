package routine

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

// LocalStorage provides goroutine-local variables.
type LocalStorage interface {
	// Get returns the value in the current goroutine's local storage, if it was set before.
	Get() (value interface{})

	// Set copy the value into the current goroutine's local storage, and return the old value.
	Set(value interface{}) (oldValue interface{})

	// Del delete the value from the current goroutine's local storage, and return it.
	Del() (oldValue interface{})
}

// Clear delete values from all goroutine's local storages.
func Clear() {
	s := loadCurrentStore(false)
	if s == nil {
		return
	}
	s.values = map[uintptr]interface{}{}
	atomic.StoreUint32(&s.count, 0)
}

// ImmutableContext represents all local storages of one goroutine.
type ImmutableContext struct {
	gid    int64
	count  uint32
	values map[uintptr]interface{}
}

// Go start an new goroutine, and copy all local storages from current goroutine.
func Go(f func()) {
	ic := BackupContext()
	go func() {
		InheritContext(ic)
		f()
	}()
}

// BackupContext copy all local storages into an ImmutableContext instance.
func BackupContext() *ImmutableContext {
	s := loadCurrentStore(false)
	if s == nil || s.values == nil || len(s.values) == 0 {
		return nil
	}
	data := make(map[uintptr]interface{}, len(s.values))
	for k, v := range s.values {
		data[k] = v
	}
	return &ImmutableContext{gid: s.gid, count: atomic.LoadUint32(&s.count), values: data}
}

// InheritContext load the specified ImmutableContext instance into the local storage of current goroutine.
func InheritContext(ic *ImmutableContext) {
	if ic == nil || ic.values == nil || len(ic.values) == 0 {
		return
	}
	s := loadCurrentStore(true)
	for k, v := range ic.values {
		s.values[k] = v
	}
	atomic.StoreUint32(&s.count, ic.count)
}

// NewLocalStorage create and return an new LocalStorage instance.
func NewLocalStorage() LocalStorage {
	t := &storage{}
	t.id = uintptr(unsafe.Pointer(t))
	return t
}

// Goid return the current goroutine's unique id.
// It will try get gid by native cgo/asm for better performance,
// and could parse gid from stack for failover supporting.
func Goid() (id int64) {
	var succ bool
	if id, succ = getGoidByNative(); !succ {
		// no need to warning
		id = getGoidByStack()
	}
	return
}

// AllGoids return all goroutine's goid in the current golang process.
// It will try load all goid from runtime natively for better performance,
// and fallover to runtime.Stack, which is realy inefficient.
func AllGoids() (ids []int64) {
	var err error
	if ids, err = getAllGoidByNative(); err != nil {
		fmt.Println("[WARNING] cannot get all goids from runtime natively, now fall over to stack info, this will be very inefficient!!!")
		ids = getAllGoidByStack()
	}
	return
}

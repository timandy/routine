package routine

import (
	"fmt"
	"reflect"
	"unsafe"
)

type g struct {
}

//go:norace
func (g *g) goid() uint64 {
	return *(*uint64)(add(unsafe.Pointer(g), offsetGoid))
}

//go:norace
func (g *g) gopc() uintptr {
	return *(*uintptr)(add(unsafe.Pointer(g), offsetGopc))
}

//go:norace
func (g *g) getPanicOnFault() bool {
	return *(*bool)(add(unsafe.Pointer(g), offsetPaniconfault))
}

//go:norace
func (g *g) setPanicOnFault(new bool) (old bool) {
	panicOnFault := (*bool)(add(unsafe.Pointer(g), offsetPaniconfault))
	old = *panicOnFault
	*panicOnFault = new
	return old
}

//go:norace
func (g *g) getLabels() unsafe.Pointer {
	return *(*unsafe.Pointer)(add(unsafe.Pointer(g), offsetLabels))
}

//go:norace
func (g *g) setLabels(labels unsafe.Pointer) {
	*(*unsafe.Pointer)(add(unsafe.Pointer(g), offsetLabels)) = labels
}

// getg returns current coroutine struct.
func getg() *g {
	gp := getgp()
	if gp == nil {
		panic("Failed to get gp from runtime natively.")
	}
	return gp
}

// offset returns the offset of the specified field.
func offset(t reflect.Type, f string) uintptr {
	field, found := t.FieldByName(f)
	if found {
		return field.Offset
	}
	panic(fmt.Sprintf("No such field '%v' of struct '%v.%v'.", f, t.PkgPath(), t.Name()))
}

// add pointer addition operation.
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

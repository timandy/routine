package routine

import (
	"fmt"
	"reflect"
	"unsafe"
)

var (
	offsetGoid         uintptr
	offsetPaniconfault uintptr
	offsetLabels       uintptr
)

func init() {
	gt := getgt()
	offsetGoid = offset(gt, "goid")
	offsetPaniconfault = offset(gt, "paniconfault")
	offsetLabels = offset(gt, "labels")
}

type g struct {
	goid         int64
	paniconfault *bool
	labels       *unsafe.Pointer
}

//go:norace
func (gp g) getPanicOnFault() bool {
	return *gp.paniconfault
}

//go:norace
func (gp g) setPanicOnFault(new bool) (old bool) {
	old = *gp.paniconfault
	*gp.paniconfault = new
	return old
}

//go:norace
func (gp g) getLabels() unsafe.Pointer {
	return *gp.labels
}

//go:norace
func (gp g) setLabels(labels unsafe.Pointer) {
	*gp.labels = labels
}

// getg returns current coroutine struct.
func getg() g {
	gp := getgp()
	if gp == nil {
		panic("Failed to get gp from runtime natively.")
	}
	return g{
		goid:         *(*int64)(add(gp, offsetGoid)),
		paniconfault: (*bool)(add(gp, offsetPaniconfault)),
		labels:       (*unsafe.Pointer)(add(gp, offsetLabels)),
	}
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

package routine

import (
	"reflect"
	"unsafe"
)

// packEface returns an empty interface representing a value of the specified type,
// using p as the pointer to the data.
//
//go:linkname packEface routine.packEface
func packEface(typ reflect.Type, p unsafe.Pointer) (i any)

// typeByString returns the type whose 'String' property equals to the given string,
// or nil if not found.
//
//go:linkname typeByString routine.typeByString
func typeByString(str string) reflect.Type

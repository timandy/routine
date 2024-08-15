package routine

import (
	"reflect"
	"unsafe"

	_ "github.com/timandy/routine/g"
)

// getgp returns the pointer to the current runtime.g.
//
//go:linkname getgp github.com/timandy/routine/g.getgp
func getgp() unsafe.Pointer

// getgt returns the type of runtime.g.
//
//go:linkname getgt github.com/timandy/routine/g.getgt
func getgt() reflect.Type

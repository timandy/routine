package routine

import (
	"reflect"
	_ "unsafe"

	_ "github.com/timandy/routine/g"
)

// getgp returns the pointer to the current runtime.g.
//
//go:linkname getgp runtime.getgp
func getgp() *g

// getgt returns the type of runtime.g.
//
//go:linkname getgt runtime.getgt
func getgt() reflect.Type

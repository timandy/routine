package routine

import (
	_ "github.com/timandy/routine/g"
	"unsafe"
)

// getg returns the pointer to the current g.
//go:linkname getg github.com/timandy/routine/g.getg
func getg() unsafe.Pointer

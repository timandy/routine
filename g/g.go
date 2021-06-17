// Copyright 2018 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

// Package g exposes goroutine struct g to user space.
package g

import (
	"unsafe"
)

func getg() unsafe.Pointer

// G returns current g (the goroutine struct) to user space.
func G() unsafe.Pointer {
	return getg()
}

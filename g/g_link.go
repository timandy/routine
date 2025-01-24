// Copyright 2021-2025 TimAndy. All rights reserved.
// Licensed under the Apache-2.0 license that can be found in the LICENSE file.

//go:build routinex

package g

import (
	"reflect"
	"unsafe"
)

// getg0 returns the value of runtime.g0.
//
//go:nosplit
//go:linkname getg0 runtime.getg0
func getg0() any

// getgp returns the pointer to the current runtime.g.
//
//go:nosplit
//go:linkname getgp runtime.getgp
func getgp() unsafe.Pointer

// getgt returns the type of runtime.g.
//
//go:nosplit
//go:linkname getgt runtime.getgt
func getgt() reflect.Type {
	return reflect.TypeOf(getg0())
}

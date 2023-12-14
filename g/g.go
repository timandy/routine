// Copyright 2021-2024 TimAndy. All rights reserved.
// Licensed under the Apache-2.0 license that can be found in the LICENSE file.

package g

import (
	"reflect"
	"unsafe"
)

// g0 the value of runtime.g0.
//
//go:linkname g0 runtime.g0
var g0 struct{}

// getgp returns the pointer to the current runtime.g.
//
//go:nosplit
func getgp() unsafe.Pointer

// getg0 returns the value of runtime.g0.
//
//go:nosplit
func getg0() interface{} {
	return packEface(getgt(), unsafe.Pointer(&g0))
}

// getgt returns the type of runtime.g.
//
//go:nosplit
func getgt() reflect.Type {
	return typeByString("runtime.g")
}

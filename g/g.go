// Copyright 2021-2024 TimAndy. All rights reserved.
// Licensed under the Apache-2.0 license that can be found in the LICENSE file.

package g

import (
	"reflect"
	"unsafe"
)

// getgp returns the pointer to the current runtime.g.
//
//go:nosplit
func getgp() unsafe.Pointer

// getgt returns the type of runtime.g.
//
//go:nosplit
func getgt() reflect.Type {
	return typeByString("runtime.g")
}

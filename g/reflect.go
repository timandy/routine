// Copyright 2021-2024 TimAndy. All rights reserved.
// Licensed under the Apache-2.0 license that can be found in the LICENSE file.

package g

import (
	"reflect"
	"unsafe"
)

// eface The empty interface struct.
type eface struct {
	_type unsafe.Pointer
	data  unsafe.Pointer
}

// iface The interface struct.
type iface struct {
	tab  unsafe.Pointer
	data unsafe.Pointer
}

// typelinks returns a slice of the sections in each module, and a slice of *rtype offsets in each module. The types in each module are sorted by string.
//
//go:linkname typelinks reflect.typelinks
func typelinks() (sections []unsafe.Pointer, offset [][]int32)

// resolveTypeOff resolves an *rtype offset from a base type.
//
//go:linkname resolveTypeOff reflect.resolveTypeOff
func resolveTypeOff(rtype unsafe.Pointer, off int32) unsafe.Pointer

// packEface returns an empty interface representing a value of the specified type, using p as the pointer to the data.
func packEface(typ reflect.Type, p unsafe.Pointer) (i interface{}) {
	t := (*iface)(unsafe.Pointer(&typ))
	e := (*eface)(unsafe.Pointer(&i))
	e._type = t.data
	e.data = p
	return
}

// typeByString returns the type whose 'String' property equals to the given string, or nil if not found.
func typeByString(str string) reflect.Type {
	// The s is search target
	s := str
	if len(str) == 0 || str[0] != '*' {
		s = "*" + s
	}
	// The typ is a struct iface{tab(ptr->reflect.Type), data(ptr->rtype)}
	typ := reflect.TypeOf(0)
	face := (*iface)(unsafe.Pointer(&typ))
	// Find the specified target through binary search algorithm
	sections, offset := typelinks()
	for offsI, offs := range offset {
		section := sections[offsI]
		// We are looking for the first index i where the string becomes >= s.
		// This is a copy of sort.Search, with f(h) replaced by (*typ[h].String() >= s).
		i, j := 0, len(offs)
		for i < j {
			h := i + (j-i)/2 // avoid overflow when computing h
			// i ≤ h < j
			face.data = resolveTypeOff(section, offs[h])
			if !(typ.String() >= s) {
				i = h + 1 // preserves f(i-1) == false
			} else {
				j = h // preserves f(j) == true
			}
		}
		// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
		// Having found the first, linear scan forward to find the last.
		// We could do a second binary search, but the caller is going
		// to do a linear scan anyway.
		if i < len(offs) {
			face.data = resolveTypeOff(section, offs[i])
			if typ.Kind() == reflect.Ptr {
				if typ.String() == str {
					return typ
				}
				elem := typ.Elem()
				if elem.String() == str {
					return elem
				}
			}
		}
	}
	return nil
}

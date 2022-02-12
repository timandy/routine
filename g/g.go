// Copyright 2018 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package g

import "unsafe"

// getg returns the pointer to the current g.
func getg() unsafe.Pointer

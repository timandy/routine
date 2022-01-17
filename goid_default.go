//go:build !go1.16
// +build !go1.16

package routine

import "errors"

var errUnsupported = errors.New("unsupported")

func getAllGoidByNative() ([]int64, bool) {
	return nil, false
}

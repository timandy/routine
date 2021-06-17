// +build !go1.16

package routine

import "errors"

var errUnsupported = errors.New("unsupported")

func getAllGoidByNative() (goids []int64, err error) {
	return nil, errUnsupported
}

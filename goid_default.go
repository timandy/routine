//go:build !go1.16
// +build !go1.16

package routine

func getAllGoidByNative() ([]int64, bool) {
	return nil, false
}

//go:build !go1.12
// +build !go1.12

package routine

func getAllGoidByNative() ([]int64, bool) {
	return nil, false
}

func foreachGoidByNative(fun func(goid int64)) bool {
	return false
}

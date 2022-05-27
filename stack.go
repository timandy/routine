package routine

import "runtime"

func captureStackTrace(skip int, depth int) []uintptr {
	pcs := make([]uintptr, depth)
	return pcs[:runtime.Callers(skip+2, pcs)]
}

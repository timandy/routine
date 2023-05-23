package routine

import (
	"runtime"
	"strings"
)

const (
	runtimePkgPrefix = "runtime."
	runtimePanic     = "panic"
)

func captureStackTrace(skip int, depth int) []uintptr {
	pcs := make([]uintptr, depth)
	return pcs[:runtime.Callers(skip+2, pcs)]
}

func showFrame(name string) bool {
	return strings.IndexByte(name, '.') >= 0 && (!strings.HasPrefix(name, runtimePkgPrefix) || isExportedRuntime(name))
}

func skipFrame(name string, skipped bool) bool {
	return !skipped && isPanicRuntime(name)
}

func isExportedRuntime(name string) bool {
	const n = len(runtimePkgPrefix)
	return len(name) > n && name[:n] == runtimePkgPrefix && 'A' <= name[n] && name[n] <= 'Z'
}

func isPanicRuntime(name string) bool {
	const n = len(runtimePkgPrefix)
	return len(name) > n && name[:n] == runtimePkgPrefix && strings.Contains(strings.ToLower(name[n:]), runtimePanic)
}

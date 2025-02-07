//go:build routinex

package routine

const routinexEnabled = true

var (
	offsetGoid         uintptr
	offsetPaniconfault uintptr
	offsetGopc         uintptr
	offsetLabels       uintptr
	offsetThreadLocals uintptr
)

func init() {
	gt := getgt()
	offsetGoid = offset(gt, "goid")
	offsetPaniconfault = offset(gt, "paniconfault")
	offsetGopc = offset(gt, "gopc")
	offsetLabels = offset(gt, "labels")
	offsetThreadLocals = offset(gt, "threadLocals")
}

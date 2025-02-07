//go:build !routinex

package routine

const routinexEnabled = false

var (
	offsetGoid         uintptr
	offsetPaniconfault uintptr
	offsetGopc         uintptr
	offsetLabels       uintptr
)

func init() {
	gt := getgt()
	offsetGoid = offset(gt, "goid")
	offsetPaniconfault = offset(gt, "paniconfault")
	offsetGopc = offset(gt, "gopc")
	offsetLabels = offset(gt, "labels")
}

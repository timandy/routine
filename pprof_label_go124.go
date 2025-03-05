//go:build go1.24

package routine

type labelMap []any

func (m labelMap) isEmpty() bool {
	return len(m) == 0
}

func defaultLabels() labelMap {
	return labelMap{}
}

//go:build !go1.24

package routine

type labelMap map[string]string

func (m labelMap) isEmpty() bool {
	return len(m) == 0
}

func defaultLabels() labelMap {
	return nil
}

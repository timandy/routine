package routine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLabelMap_IsEmpty(t *testing.T) {
	labels := labelMap{}
	assert.True(t, labels.isEmpty())
}

func TestDefaultLabels(t *testing.T) {
	labels := defaultLabels()
	assert.True(t, labels.isEmpty())
}

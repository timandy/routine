package routine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLabelMap_Nil(t *testing.T) {
	var labels labelMap
	assert.Nil(t, labels)
	//
	labels = labelMap{}
	assert.NotNil(t, labels)
}

func TestLabelMap_Empty(t *testing.T) {
	var labels labelMap
	assert.Empty(t, labels)
	//
	labels = labelMap{}
	assert.Empty(t, labels)
}

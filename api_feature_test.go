package routine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFeature(t *testing.T) {
	fea := NewFeature()
	assert.NotNil(t, fea)
	//
	p, ok := fea.(*feature)
	assert.Same(t, p, fea)
	assert.True(t, ok)
}

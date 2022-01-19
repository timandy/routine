package routine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFeature(t *testing.T) {
	fea := NewFeature()
	assert.NotNil(t, fea)
	//
	p, ok := fea.(*feature)
	assert.Same(t, p, fea)
	assert.True(t, ok)
}

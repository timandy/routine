package routine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFuture(t *testing.T) {
	fea := NewFuture()
	assert.NotNil(t, fea)
	//
	p, ok := fea.(*future)
	assert.Same(t, p, fea)
	assert.True(t, ok)
}

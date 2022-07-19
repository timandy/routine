package routine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFuture(t *testing.T) {
	fut := NewFuture()
	assert.NotNil(t, fut)
	//
	p, ok := fut.(*future)
	assert.Same(t, p, fut)
	assert.True(t, ok)
}

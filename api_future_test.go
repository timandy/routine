package routine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFuture(t *testing.T) {
	fut := NewFuture()
	assert.NotNil(t, fut)
	//
	p, ok := fut.(*future)
	assert.Same(t, p, fut)
	assert.True(t, ok)
}

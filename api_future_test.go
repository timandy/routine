package routine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCancelToken(t *testing.T) {
	fut := NewFuture[any]()
	token, ok := fut.(CancelToken)
	assert.Same(t, fut, token)
	assert.True(t, ok)
}

func TestNewFuture(t *testing.T) {
	fut := NewFuture[any]()
	assert.NotNil(t, fut)
	//
	p, ok := fut.(*future[any])
	assert.Same(t, p, fut)
	assert.True(t, ok)
}

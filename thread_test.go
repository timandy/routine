package routine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCurrentThread(t *testing.T) {
	assert.Nil(t, currentThread(false))
	assert.NotNil(t, currentThread(true))
	assert.Same(t, currentThread(false), currentThread(true))
}

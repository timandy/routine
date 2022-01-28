package routine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHash(t *testing.T) {
	assert.Equal(t, 1, hash(-17))
	assert.Equal(t, 0, hash(-16))
	assert.Equal(t, 15, hash(-15))
	assert.Equal(t, 1, hash(-1))
	assert.Equal(t, 0, hash(0))
	assert.Equal(t, 1, hash(1))
	assert.Equal(t, 15, hash(15))
	assert.Equal(t, 0, hash(16))
	assert.Equal(t, 1, hash(17))
}

func TestCurrentThread(t *testing.T) {
	assert.Nil(t, currentThread(false))
	assert.NotNil(t, currentThread(true))
	assert.Same(t, currentThread(false), currentThread(true))
}

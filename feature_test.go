package routine

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestComplete(t *testing.T) {
	fea := NewFeature()
	go func() {
		fea.Complete(1)
	}()
	assert.Equal(t, 1, fea.Get())
}

func TestCompleteError(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			re := err.(RuntimeError)
			assert.NotNil(t, re)
			assert.Equal(t, "1", re.Message())
			assert.NotNil(t, re.StackTrace())
		}
	}()

	fea := NewFeature()
	go func() {
		fea.CompleteError(1)
	}()
	fea.Get()
	assert.Fail(t, "should not be here")
}

func TestGet(t *testing.T) {
	run := false
	fea := NewFeature()
	go func() {
		time.Sleep(500 * time.Millisecond)
		run = true
		fea.Complete(nil)
	}()
	assert.Nil(t, fea.Get())
	assert.True(t, run)
}

package routine

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFuture_Complete(t *testing.T) {
	fea := NewFeature()
	go func() {
		fea.Complete(1)
	}()
	assert.Equal(t, 1, fea.Get())
}

func TestFuture_CompleteError(t *testing.T) {
	defer func() {
		if cause := recover(); cause != nil {
			err := cause.(RuntimeError)
			assert.NotNil(t, err)
			assert.Equal(t, "1", err.Message())
			assert.NotNil(t, err.StackTrace())
		}
	}()

	fea := NewFeature()
	go func() {
		fea.CompleteError(1)
	}()
	fea.Get()
	assert.Fail(t, "should not be here")
}

func TestFuture_Get(t *testing.T) {
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

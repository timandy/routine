package routine

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	feature := NewFeature()
	begin := time.Now()
	go func() {
		time.Sleep(time.Millisecond * 500)
		feature.Complete(nil)
	}()
	assert.Nil(t, feature.Get())
	assert.Greater(t, time.Now().Sub(begin), time.Millisecond*200)
}

func TestComplete(t *testing.T) {
	feature := NewFeature()
	go func() {
		feature.Complete(1)
	}()
	assert.Equal(t, 1, feature.Get())
}

func TestCompleteError(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			stackError := err.(*StackError)
			assert.NotNil(t, stackError)
			assert.Equal(t, 1, stackError.error)
			assert.NotNil(t, stackError.stackTrace)
		}
	}()

	feature := NewFeature()
	go func() {
		feature.CompleteError(1)
	}()
	feature.Get()
	assert.Fail(t, "should not be here")
}

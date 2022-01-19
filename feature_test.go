package routine

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	fea := NewFeature()
	begin := time.Now()
	go func() {
		time.Sleep(time.Millisecond * 500)
		fea.Complete(nil)
	}()
	assert.Nil(t, fea.Get())
	assert.Greater(t, time.Now().Sub(begin), time.Millisecond*200)
}

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
			stackError := err.(*StackError)
			assert.NotNil(t, stackError)
			assert.Equal(t, 1, stackError.error)
			assert.NotNil(t, stackError.stackTrace)
		}
	}()

	fea := NewFeature()
	go func() {
		fea.CompleteError(1)
	}()
	fea.Get()
	assert.Fail(t, "should not be here")
}

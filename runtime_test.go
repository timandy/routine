package routine

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func TestGetgLink(t *testing.T) {
	gp := getg()
	assert.NotNil(t, gp)
}

func TestCurGoroutineID(t *testing.T) {
	assert.NotEqual(t, 0, curGoroutineID())
}

func TestProfLabel(t *testing.T) {
	value := "Hello World"
	setProfLabel(unsafe.Pointer(&value))
	//
	GoWait(func() {
		assert.Equal(t, value, *(*string)(getProfLabel()))
		//
		value2 := "Hello 世界"
		setProfLabel(unsafe.Pointer(&value2))
		assert.Equal(t, value2, *(*string)(getProfLabel()))
	}).Get()
	//
	assert.Equal(t, value, *(*string)(getProfLabel()))
}

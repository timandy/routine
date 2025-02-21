package routine

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestPackEface(t *testing.T) {
	value := 1
	valueInterface := packEface(reflect.TypeOf(0), unsafe.Pointer(&value))
	assert.Equal(t, value, valueInterface)
	//
	value = 2
	assert.Equal(t, value, valueInterface)
}

func TestTypeByString(t *testing.T) {
	gt := typeByString("runtime.g")
	assert.NotNil(t, gt)
	assert.Equal(t, "runtime.g", gt.String())
	fGoid, ok := gt.FieldByName("goid")
	assert.True(t, ok)
	assert.Greater(t, int(fGoid.Offset), 0)
	//
	gt2 := typeByString("*runtime.g")
	assert.NotNil(t, gt2)
	assert.Equal(t, "*runtime.g", gt2.String())
	fGoid2, ok2 := gt2.Elem().FieldByName("goid")
	assert.True(t, ok2)
	assert.Greater(t, int(fGoid2.Offset), 0)
	assert.Equal(t, fGoid.Offset, fGoid2.Offset)
	//
	assert.Nil(t, typeByString("runtime.Pointer"))
}

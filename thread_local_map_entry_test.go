package routine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntry_Clone(t *testing.T) {
	expect := &personCloneable{Id: 1, Name: "Hello"}
	e := entry(expect)
	value := entryValue[*personCloneable](e)
	assert.NotNil(t, value)
	assert.Equal(t, *expect, *value)
	//
	c, ok := entryAssert[Cloneable](e)
	assert.True(t, ok)
	assert.Equal(t, *expect, *c.(*personCloneable))
	//
	copied := entry(c.Clone())
	value2 := entryValue[*personCloneable](copied)
	assert.NotNil(t, value2)
	assert.Equal(t, *expect, *value2)
	//
	c3, ok2 := entryAssert[Cloneable](copied)
	assert.True(t, ok2)
	assert.Equal(t, *expect, *c3.(*personCloneable))
}

//===

func TestEntry_Value_Nil(t *testing.T) {
	e := entry(nil)
	value := entryValue[any](e)
	assert.Nil(t, value)
}

func TestEntry_Value_NotNil(t *testing.T) {
	expect := 1
	e := entry(expect)
	value := entryValue[int](e)
	assert.Equal(t, expect, value)
}

func TestEntry_Value_Default(t *testing.T) {
	expect := 0
	e := entry(expect)
	value := entryValue[int](e)
	assert.Equal(t, expect, value)
}

func TestEntry_Value_Interface_Nil(t *testing.T) {
	var expect Cloneable
	e := entry(expect)
	value := entryValue[Cloneable](e)
	assert.Nil(t, value)
}

func TestEntry_Value_Interface_NotNil(t *testing.T) {
	var expect Cloneable = &personCloneable{Id: 1, Name: "Hello"}
	e := entry(expect)
	value := entryValue[Cloneable](e)
	assert.Same(t, expect, value)
}

func TestEntry_Value_Pointer_Nil(t *testing.T) {
	var expect *personCloneable
	e := entry(expect)
	value := entryValue[*personCloneable](e)
	assert.Nil(t, value)
}

func TestEntry_Value_Pointer_NotNil(t *testing.T) {
	expect := &personCloneable{Id: 1, Name: "Hello"}
	e := entry(expect)
	value := entryValue[*personCloneable](e)
	assert.Same(t, expect, value)
}

func TestEntry_Value_Struct_Default(t *testing.T) {
	expect := personCloneable{}
	e := entry(expect)
	value := entryValue[personCloneable](e)
	assert.Equal(t, expect, value)
}

func TestEntry_Value_Struct_NotDefault(t *testing.T) {
	expect := personCloneable{Id: 1, Name: "Hello"}
	e := entry(expect)
	value := entryValue[personCloneable](e)
	assert.Equal(t, expect, value)
}

//===

func TestEntry_Assert_Nil(t *testing.T) {
	e := entry(nil)
	value, ok := entryAssert[any](e)
	assert.False(t, ok)
	assert.Nil(t, value)
}

func TestEntry_Assert_NotNil(t *testing.T) {
	expect := 1
	e := entry(expect)
	value, ok := entryAssert[int](e)
	assert.True(t, ok)
	assert.Equal(t, expect, value)
}

func TestEntry_Assert_Default(t *testing.T) {
	expect := 0
	e := entry(expect)
	value, ok := entryAssert[int](e)
	assert.True(t, ok)
	assert.Equal(t, expect, value)
}

func TestEntry_Assert_Interface_Nil(t *testing.T) {
	var expect Cloneable
	e := entry(expect)
	value, ok := entryAssert[Cloneable](e)
	assert.False(t, ok)
	assert.Nil(t, value)
}

func TestEntry_Assert_Interface_NotNil(t *testing.T) {
	var expect Cloneable = &personCloneable{Id: 1, Name: "Hello"}
	e := entry(expect)
	value, ok := entryAssert[Cloneable](e)
	assert.True(t, ok)
	assert.Same(t, expect, value)
}

func TestEntry_Assert_Pointer_Nil(t *testing.T) {
	var expect *personCloneable
	e := entry(expect)
	value, ok := entryAssert[*personCloneable](e)
	assert.True(t, ok)
	assert.Nil(t, value)
}

func TestEntry_Assert_Pointer_NotNil(t *testing.T) {
	expect := &personCloneable{Id: 1, Name: "Hello"}
	e := entry(expect)
	value, ok := entryAssert[*personCloneable](e)
	assert.True(t, ok)
	assert.Same(t, expect, value)
}

func TestEntry_Assert_Struct_Default(t *testing.T) {
	expect := personCloneable{}
	e := entry(expect)
	value, ok := entryAssert[personCloneable](e)
	assert.True(t, ok)
	assert.Equal(t, expect, value)
}

func TestEntry_Assert_Struct_NotDefault(t *testing.T) {
	expect := personCloneable{Id: 1, Name: "Hello"}
	e := entry(expect)
	value, ok := entryAssert[personCloneable](e)
	assert.True(t, ok)
	assert.Equal(t, expect, value)
}

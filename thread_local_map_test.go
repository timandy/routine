package routine

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestCreateInheritedMapValue(t *testing.T) {
	tls := NewInheritableThreadLocal()
	value := rand.Uint64()
	tls.Set(value)
	srcValue := tls.Get()
	assert.NotSame(t, &value, &srcValue)
	assert.Equal(t, value, srcValue)

	mp := createInheritedMap(currentThread().inheritableThreadLocals)
	assert.NotNil(t, mp)
	getValue := mp.get(tls)
	assert.NotSame(t, &value, &getValue)
	assert.Equal(t, value, getValue)

	mp2 := createInheritedMap(currentThread().inheritableThreadLocals)
	assert.NotNil(t, mp2)
	assert.NotSame(t, mp, mp2)
	getValue2 := mp2.get(tls)
	assert.NotSame(t, &value, &getValue2)
	assert.Equal(t, value, getValue2)
}

func TestCreateInheritedMapStruct(t *testing.T) {
	tls := NewInheritableThreadLocal()
	value := personCloneable{Id: 1, Name: "Hello"}
	tls.Set(value)
	srcValue := tls.Get()
	assert.NotSame(t, &value, &srcValue)
	assert.Equal(t, value, srcValue)

	mp := createInheritedMap(currentThread().inheritableThreadLocals)
	assert.NotNil(t, mp)
	getValue := mp.get(tls)
	assert.NotSame(t, &value, &getValue)
	assert.Equal(t, value, getValue)

	mp2 := createInheritedMap(currentThread().inheritableThreadLocals)
	assert.NotNil(t, mp2)
	assert.NotSame(t, mp, mp2)
	getValue2 := mp2.get(tls)
	assert.NotSame(t, &value, &getValue2)
	assert.Equal(t, value, getValue2)
}

func TestCreateInheritedMapPointer(t *testing.T) {
	tls := NewInheritableThreadLocal()
	value := &person{Id: 1, Name: "Hello"}
	tls.Set(value)
	srcValue := tls.Get().(*person)
	assert.Same(t, value, srcValue)
	assert.Equal(t, *value, *srcValue)

	mp := createInheritedMap(currentThread().inheritableThreadLocals)
	assert.NotNil(t, mp)
	getValue := mp.get(tls).(*person)
	assert.Same(t, value, getValue)
	assert.Equal(t, *value, *getValue)

	mp2 := createInheritedMap(currentThread().inheritableThreadLocals)
	assert.NotNil(t, mp2)
	assert.NotSame(t, mp, mp2)
	getValue2 := mp2.get(tls).(*person)
	assert.Same(t, value, getValue2)
	assert.Equal(t, *value, *getValue2)
}

func TestCreateInheritedMapCloneable(t *testing.T) {
	tls := NewInheritableThreadLocal()
	value := &personCloneable{Id: 1, Name: "Hello"}
	tls.Set(value)
	srcValue := tls.Get().(*personCloneable)
	assert.Same(t, value, srcValue)
	assert.Equal(t, *value, *srcValue)

	mp := createInheritedMap(currentThread().inheritableThreadLocals)
	assert.NotNil(t, mp)
	getValue := mp.get(tls).(*personCloneable)
	assert.NotSame(t, value, getValue)
	assert.Equal(t, *value, *getValue)

	mp2 := createInheritedMap(currentThread().inheritableThreadLocals)
	assert.NotNil(t, mp2)
	assert.NotSame(t, mp, mp2)
	getValue2 := mp2.get(tls).(*personCloneable)
	assert.NotSame(t, value, getValue2)
	assert.Equal(t, *value, *getValue2)
}

func TestFill(t *testing.T) {
	a := make([]Any, 6)
	fill(a, 4, 5, unset)
	for i := 0; i < 6; i++ {
		if i == 4 {
			assert.True(t, a[i] == unset)
		} else {
			assert.Nil(t, a[i])
			assert.True(t, a[i] != unset)
		}
	}
}

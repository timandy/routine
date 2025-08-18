package routine

import (
	"math/rand"
	"sync"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

//go:linkname createInheritedMapExport routine.createInheritedMap
func createInheritedMapExport() unsafe.Pointer

//go:linkname restoreInheritedMapExport routine.restoreInheritedMap
func restoreInheritedMapExport(mp unsafe.Pointer) func()

func TestObject(t *testing.T) {
	var value entry = &object{}
	assert.NotSame(t, unset, value)
	//
	var value2 entry = &object{}
	assert.NotSame(t, value, value2)
	//
	var value3 any = unset
	assert.Same(t, unset, value3)
}

func TestCreateInheritedMap(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		thd := currentThread(true)
		assert.NotNil(t, thd)
		assert.Nil(t, thd.inheritableThreadLocals)
		thd.inheritableThreadLocals = &threadLocalMap{}
		assert.Nil(t, thd.inheritableThreadLocals.table)
		assert.Nil(t, createInheritedMap())
		//
		wg.Done()
	}()
	wg.Wait()
}

func TestCreateInheritedMap_Nil(t *testing.T) {
	tls := NewInheritableThreadLocal[string]()
	tls.Set("")
	srcValue := tls.Get()
	assert.Equal(t, "", srcValue)
	assert.True(t, srcValue == "")

	mp := createInheritedMap()
	assert.NotNil(t, mp)
	getValue := entryValue[string](mp.get(tls.(*inheritableThreadLocal[string]).index))
	assert.Equal(t, "", getValue)
	assert.True(t, getValue == "")

	mp2 := createInheritedMap()
	assert.NotNil(t, mp2)
	assert.NotSame(t, mp, mp2)
	getValue2 := entryValue[string](mp2.get(tls.(*inheritableThreadLocal[string]).index))
	assert.Equal(t, "", getValue2)
	assert.True(t, getValue2 == "")
}

func TestCreateInheritedMap_Value(t *testing.T) {
	tls := NewInheritableThreadLocal[uint64]()
	value := rand.Uint64()
	tls.Set(value)
	srcValue := tls.Get()
	assert.NotSame(t, &value, &srcValue)
	assert.Equal(t, value, srcValue)

	mp := createInheritedMap()
	assert.NotNil(t, mp)
	getValue := entryValue[uint64](mp.get(tls.(*inheritableThreadLocal[uint64]).index))
	assert.NotSame(t, &value, &getValue)
	assert.Equal(t, value, getValue)

	mp2 := createInheritedMap()
	assert.NotNil(t, mp2)
	assert.NotSame(t, mp, mp2)
	getValue2 := entryValue[uint64](mp2.get(tls.(*inheritableThreadLocal[uint64]).index))
	assert.NotSame(t, &value, &getValue2)
	assert.Equal(t, value, getValue2)
}

func TestCreateInheritedMap_Struct(t *testing.T) {
	tls := NewInheritableThreadLocal[personCloneable]()
	value := personCloneable{Id: 1, Name: "Hello"}
	tls.Set(value)
	srcValue := tls.Get()
	assert.NotSame(t, &value, &srcValue)
	assert.Equal(t, value, srcValue)

	mp := createInheritedMap()
	assert.NotNil(t, mp)
	getValue := entryValue[personCloneable](mp.get(tls.(*inheritableThreadLocal[personCloneable]).index))
	assert.NotSame(t, &value, &getValue)
	assert.Equal(t, value, getValue)

	mp2 := createInheritedMap()
	assert.NotNil(t, mp2)
	assert.NotSame(t, mp, mp2)
	getValue2 := entryValue[personCloneable](mp2.get(tls.(*inheritableThreadLocal[personCloneable]).index))
	assert.NotSame(t, &value, &getValue2)
	assert.Equal(t, value, getValue2)
}

func TestCreateInheritedMap_Pointer(t *testing.T) {
	tls := NewInheritableThreadLocal[*person]()
	value := &person{Id: 1, Name: "Hello"}
	tls.Set(value)
	srcValue := tls.Get()
	assert.Same(t, value, srcValue)
	assert.Equal(t, *value, *srcValue)

	mp := createInheritedMap()
	assert.NotNil(t, mp)
	getValue := entryValue[*person](mp.get(tls.(*inheritableThreadLocal[*person]).index))
	assert.Same(t, value, getValue)
	assert.Equal(t, *value, *getValue)

	mp2 := createInheritedMap()
	assert.NotNil(t, mp2)
	assert.NotSame(t, mp, mp2)
	getValue2 := entryValue[*person](mp2.get(tls.(*inheritableThreadLocal[*person]).index))
	assert.Same(t, value, getValue2)
	assert.Equal(t, *value, *getValue2)
}

func TestCreateInheritedMap_Cloneable(t *testing.T) {
	tls := NewInheritableThreadLocal[*personCloneable]()
	value := &personCloneable{Id: 1, Name: "Hello"}
	tls.Set(value)
	srcValue := tls.Get()
	assert.Same(t, value, srcValue)
	assert.Equal(t, *value, *srcValue)

	mp := createInheritedMap()
	assert.NotNil(t, mp)
	getValue := entryValue[*personCloneable](mp.get(tls.(*inheritableThreadLocal[*personCloneable]).index))
	assert.NotSame(t, value, getValue)
	assert.Equal(t, *value, *getValue)

	mp2 := createInheritedMap()
	assert.NotNil(t, mp2)
	assert.NotSame(t, mp, mp2)
	getValue2 := entryValue[*personCloneable](mp2.get(tls.(*inheritableThreadLocal[*personCloneable]).index))
	assert.NotSame(t, value, getValue2)
	assert.Equal(t, *value, *getValue2)
}

func TestRestoreInheritedMap(t *testing.T) {
	tls := NewInheritableThreadLocal[*personCloneable]()
	value := &personCloneable{Id: 1, Name: "Hello"}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		mp := createInheritedMap()
		assert.Nil(t, mp)
		//
		go func() {
			defer func() {
				if routinexEnabled {
					assert.NotNil(t, currentThread(false))
				} else {
					assert.Nil(t, currentThread(false))
				}
				assert.Nil(t, tls.Get())
				wg.Done()
			}()
			defer restoreInheritedMap(mp)()
		}()
		wg.Done()
	}()
	wg.Wait()
	//
	wg2 := sync.WaitGroup{}
	wg2.Add(2)
	go func() {
		mp := createInheritedMap()
		assert.Nil(t, mp)
		//
		go func() {
			defer func() {
				assert.NotNil(t, currentThread(false))
				assert.Nil(t, tls.Get())
				wg2.Done()
			}()
			defer restoreInheritedMap(mp)()
			tls.Set(value)
			assert.Same(t, value, tls.Get())
		}()
		wg2.Done()
	}()
	wg2.Wait()
	//
	value2 := &personCloneable{Id: 2, Name: "World"}
	wg3 := sync.WaitGroup{}
	wg3.Add(2)
	go func() {
		tls.Set(value)
		assert.Same(t, value, tls.Get())
		mp := createInheritedMap()
		assert.NotNil(t, mp)
		//
		go func() {
			defer func() {
				assert.NotNil(t, currentThread(false))
				assert.Same(t, value2, tls.Get())
				wg3.Done()
			}()
			tls.Set(value2)
			assert.Same(t, value2, tls.Get())
			defer restoreInheritedMap(mp)()
			result := tls.Get()
			assert.NotNil(t, result)
			assert.NotSame(t, value, result)
			assert.NotSame(t, value2, result)
			assert.Equal(t, *value, *result)
		}()
		wg3.Done()
	}()
	wg3.Wait()
}

func TestRestoreInheritedMap_Export(t *testing.T) {
	tls := NewInheritableThreadLocal[*personCloneable]()
	value := &personCloneable{Id: 1, Name: "Hello"}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		mp := createInheritedMapExport()
		assert.Nil(t, mp)
		//
		go func() {
			defer func() {
				if routinexEnabled {
					assert.NotNil(t, currentThread(false))
				} else {
					assert.Nil(t, currentThread(false))
				}
				assert.Nil(t, tls.Get())
				wg.Done()
			}()
			defer restoreInheritedMapExport(mp)()
		}()
		wg.Done()
	}()
	wg.Wait()
	//
	wg2 := sync.WaitGroup{}
	wg2.Add(2)
	go func() {
		mp := createInheritedMapExport()
		assert.Nil(t, mp)
		//
		go func() {
			defer func() {
				assert.NotNil(t, currentThread(false))
				assert.Nil(t, tls.Get())
				wg2.Done()
			}()
			defer restoreInheritedMapExport(mp)()
			tls.Set(value)
			assert.Same(t, value, tls.Get())
		}()
		wg2.Done()
	}()
	wg2.Wait()
	//
	value2 := &personCloneable{Id: 2, Name: "World"}
	wg3 := sync.WaitGroup{}
	wg3.Add(2)
	go func() {
		tls.Set(value)
		assert.Same(t, value, tls.Get())
		mp := createInheritedMapExport()
		assert.NotNil(t, mp)
		//
		go func() {
			defer func() {
				assert.NotNil(t, currentThread(false))
				assert.Same(t, value2, tls.Get())
				wg3.Done()
			}()
			tls.Set(value2)
			assert.Same(t, value2, tls.Get())
			defer restoreInheritedMapExport(mp)()
			result := tls.Get()
			assert.NotNil(t, result)
			assert.NotSame(t, value, result)
			assert.NotSame(t, value2, result)
			assert.Equal(t, *value, *result)
		}()
		wg3.Done()
	}()
	wg3.Wait()
}

func TestFill(t *testing.T) {
	a := make([]entry, 6)
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

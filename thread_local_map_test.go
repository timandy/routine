package routine

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestCreateInheritedMap(t *testing.T) {
	tls := NewInheritableThreadLocal()
	value := rand.Uint64()
	tls.Set(value)
	assert.Equal(t, value, tls.Get())

	mp := createInheritedMap()
	assert.NotNil(t, mp)
	assert.Equal(t, value, mp.get(tls))

	mp2 := createInheritedMap()
	assert.NotNil(t, mp2)
	assert.NotSame(t, mp, mp2)
	assert.Equal(t, value, mp2.get(tls))
}

func TestFill(t *testing.T) {
	a := make([]Any, 6)
	fill(a, 4, 5, 1)
	for i := 0; i < 6; i++ {
		if i == 4 {
			assert.Equal(t, 1, a[i])
		} else {
			assert.Nil(t, a[i])
		}
	}
}

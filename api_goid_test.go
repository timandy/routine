package routine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoid(t *testing.T) {
	assert.NotEqual(t, 0, Goid())
	assert.Equal(t, Goid(), Goid())
}

//===

// BenchmarkGoid-4                                229975064                4.939 ns/op            0 B/op          0 allocs/op
func BenchmarkGoid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Goid()
	}
}

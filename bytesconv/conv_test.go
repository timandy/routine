package bytesconv

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"unsafe"
)

func getBytesDataPointer(value []byte) unsafe.Pointer {
	return unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&value)).Data)
}

func getStringDataPointer(value string) unsafe.Pointer {
	return unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&value)).Data)
}

func TestBytesSlow(t *testing.T) {
	str := "Hello 世界"
	buf := []byte(str)
	assert.NotEqual(t, getStringDataPointer(str), getBytesDataPointer(buf))
}

func TestStringSlow(t *testing.T) {
	buf := []byte("Hello 世界")
	str := string(buf)
	assert.NotEqual(t, getStringDataPointer(str), getBytesDataPointer(buf))
}

func TestBytes(t *testing.T) {
	str := "Hello 世界"
	buf := Bytes(str)
	assert.Equal(t, buf, []byte(str))
	assert.Equal(t, getStringDataPointer(str), getBytesDataPointer(buf))
}

func TestString(t *testing.T) {
	buf := []byte("Hello 世界")
	str := String(buf)
	assert.Equal(t, str, string(buf))
	assert.Equal(t, getStringDataPointer(str), getBytesDataPointer(buf))
}

//===

// BenchmarkBytesSlow-4             237312201               4.475 ns/op           0 B/op          0 allocs/op
func BenchmarkBytesSlow(b *testing.B) {
	str := "Hello 世界"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = []byte(str)
	}
}

// BenchmarkStringSlow-4            282549794               3.999 ns/op           0 B/op          0 allocs/op
func BenchmarkStringSlow(b *testing.B) {
	buf := []byte("Hello 世界")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = string(buf)
	}
}

// BenchmarkBytes-4                1000000000               0.2533 ns/op          0 B/op          0 allocs/op
func BenchmarkBytes(b *testing.B) {
	str := "Hello 世界"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Bytes(str)
	}
}

// BenchmarkString-4               1000000000               0.2560 ns/op          0 B/op          0 allocs/op
func BenchmarkString(b *testing.B) {
	buf := []byte("Hello 世界")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = String(buf)
	}
}

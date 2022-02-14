package bytesconv

import (
	"reflect"
	"unsafe"
)

//Bytes convert string to bytes using zero-copy.
func Bytes(value string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&value))
	sliceHeader := &reflect.SliceHeader{
		Data: stringHeader.Data,
		Cap:  stringHeader.Len,
		Len:  stringHeader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(sliceHeader))
}

//String convert bytes to string using zero-copy.
func String(value []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&value))
	stringHeader := &reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}
	return *(*string)(unsafe.Pointer(stringHeader))
}

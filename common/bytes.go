package common

import (
	"unsafe"
)

// Copy Unsafe Ptr Given An Unsafe Ptr and Bytes Size into a Go Byte Slice
func GoBytes(ptr uintptr, size int) []byte {
	buffer := make([]byte, size)
	for i := 0; i < size; i++ {
		buffer[i] = *(*byte)(unsafe.Pointer(uintptr(ptr) + uintptr(i)))
	}
	return buffer
}

func ReflectBytes(ptr uintptr, size int) []byte {
	slice := unsafe.Slice((*byte)(unsafe.Pointer(ptr)), size)
	if slice == nil {
		panic("reflect.SliceHeader: slice is nil")
	}
	return slice
}

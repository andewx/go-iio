package common

import (
	"unsafe"
)

// GoBytes Converts an unsafe pointer to a byte slice
func GoBytes(ptr uintptr, size int) []byte {
	buffer := make([]byte, size)
	for i := 0; i < size; i++ {
		buffer[i] = *(*byte)(unsafe.Pointer(uintptr(ptr) + uintptr(i)))
	}
	return buffer
}

// ReflectBytes Converts an unsafe pointer to a byte slice
func ReflectBytes(ptr uintptr, size int) []byte {
	slice := unsafe.Slice((*byte)(unsafe.Pointer(ptr)), size)
	if slice == nil {
		panic("reflect.SliceHeader: slice is nil")
	}
	return slice
}

// ConvertByteInt16ToFloat64 Converts a byte slice to a float64 slice
func ConvertByteInt16ToFloat32(data []byte) []float32 {

	//Cast byte slice to int16 slice
	int16Data := *(*[]int16)(unsafe.Pointer(&data))

	float64Data := make([]float32, len(int16Data))

	for i := 0; i < len(int16Data); i++ {
		float64Data[i] = float32(int16Data[i])
	}

	return float64Data
}

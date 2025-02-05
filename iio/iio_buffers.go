package iio

// #cgo pkg-config: libiio
// #include <iio.h>
import "C"

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/andewx/go-iio/common"
)

// Buffer wraps an iio_buffer, which allows for reading samples from and writing
// data to the underlying device.
type Buffer struct {
	closed *bool
	handle *C.struct_iio_buffer
}

// Close will destroy the handle to the iio_buffer.
func (b Buffer) Close() error {
	if *b.closed {
		return nil
	}
	C.iio_buffer_destroy(b.handle)
	*b.closed = true
	return nil
}

// Step will return the iio_buffer_step, which is the size of each coherent
// sample.
func (b Buffer) Step() uintptr {
	return uintptr(C.iio_buffer_step(b.handle))
}

// TransformUnsafeToBuffer will copy the data from the unsafe pointer to the buffer
func (b Buffer) TransformUnsafeToBuffer(chn Channel, ptr unsafe.Pointer, size int) (int, error) {
	base := uintptr(C.iio_buffer_first(b.handle, chn.handle))
	end := uintptr(C.iio_buffer_end(b.handle))
	totalBytes := int(end - base)
	if totalBytes < 0 {
		return 0, fmt.Errorf("iio: internal error during Buffer.CopyToUnsafe")
	}

	if totalBytes == 0 {
		return 0, nil
	}

	bufMemory := common.ReflectBytes(base, totalBytes)
	targetMemory := common.ReflectBytes(uintptr(ptr), size)

	i := copy(bufMemory, targetMemory)
	return i, nil
}

// TransformBufferToUnsafe will copy the data from the buffer to the unsafe pointer
func (b Buffer) TransformBufferToUnsafe(chn Channel, ptr unsafe.Pointer, size int) (int, error) {
	base := uintptr(C.iio_buffer_first(b.handle, chn.handle))
	end := uintptr(C.iio_buffer_end(b.handle))
	totalBytes := int(end - base)
	if totalBytes < 0 {
		return 0, fmt.Errorf("iio: internal error during Buffer.CopyToUnsafe")
	}

	if totalBytes == 0 {
		return 0, nil
	}

	bufMemory := common.ReflectBytes(base, totalBytes)
	targetMemory := common.ReflectBytes(uintptr(ptr), size)

	i := copy(targetMemory, bufMemory)
	return i, nil
}

// PushPartial will push the data written to the Buffer (from start to end) to the
// Device.
func (b Buffer) PushPartial(length int) (int, error) {
	i := C.iio_buffer_push_partial(b.handle, C.size_t(length))
	if i < 0 {
		return 0, syscall.Errno(-i)
	}
	return int(i), nil
}

// Push will push the data written to the Buffer (from start to end) to the
// Device.
func (b Buffer) Push() (int, error) {
	i := C.iio_buffer_push(b.handle)
	if i < 0 {
		return 0, syscall.Errno(-i)
	}
	return int(i), nil
}

// Refill will fill the Buffer up with samples from the backing device.
func (b Buffer) Refill() (int, error) {
	i := C.iio_buffer_refill(b.handle)
	if i < 0 {
		return 0, syscall.Errno(-i)
	}
	return int(i), nil
}

func (d Device) createBuffer(samplesCount int, cyclic bool) (*Buffer, error) {
	buf, err := C.iio_device_create_buffer(
		d.handle,
		C.size_t(samplesCount),
		C.bool(cyclic),
	)
	if buf == nil {
		return nil, err
	}
	var closed bool
	return &Buffer{
		handle: buf,
		closed: &closed,
	}, nil
}

// CreateBuffer will create an iio_buffer from a given device.
//
// At least one channel must be Enabled prior to this call.
func (d Device) CreateBuffer(samplesCount int) (*Buffer, error) {
	return d.createBuffer(samplesCount, false)
}

// CreateCyclicBuffer will create an iio_buffer from a given device.
//
// At least one channel must be Enabled prior to this call.
func (d Device) CreateCyclicBuffer(samplesCount int) (*Buffer, error) {
	return d.createBuffer(samplesCount, true)
}

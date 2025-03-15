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

// DataType represents the type of data values samples can take
type DataType int

const (
	DataTypeFloat32 DataType = iota
	DataTypeFloat64
	DataTypeInt8
	DataTypeInt16
	DataTypeInt32
	DataTypeInt64
	DataTypeComplex64
	DataTypeComplex128
)

// Buffer wraps an iio_buffer, which allows for reading samples from and writing
// data to the underlying device.
type Buffer struct {
	closed    *bool
	size      int
	data_type DataType
	handle    *C.struct_iio_buffer
	data      []byte
}

// Close will destroy the handle to the iio_buffer.
func (b *Buffer) Close() error {
	if *b.closed {
		return nil
	}
	C.iio_buffer_destroy(b.handle)
	b.data = nil
	*b.closed = true
	return nil
}

// Step will return the iio_buffer_step, which is the size of each coherent
// sample.
func (b *Buffer) Step() uintptr {
	return uintptr(C.iio_buffer_step(b.handle))
}

// TransformUnsafeToBuffer will copy the data from the unsafe pointer to the buffer
func (b *Buffer) Read(chn *Channel, ptr unsafe.Pointer, size int) (int, error) {
	base := uintptr(C.iio_buffer_first(b.handle, chn.handle))
	end := uintptr(C.iio_buffer_end(b.handle))
	totalBytes := int(end - base)
	if totalBytes < 0 {
		return 0, fmt.Errorf("iio: internal error during Buffer.CopyToUnsafe")
	}

	if totalBytes == 0 {
		return 0, nil
	}

	b.data = common.ReflectBytes(base, totalBytes)

	return totalBytes, nil
}

// TransformBufferToUnsafe will copy the data from the buffer to the unsafe pointer
func (b *Buffer) Write(ptr unsafe.Pointer, size int) {

	if b.size <= 0 {
		fmt.Printf("iio: buffer size is not set")
		return
	}

	if size > len(b.data) {
		fmt.Printf("iio: size is greater than the buffer size")
		return
	}

	if size <= 0 {
		fmt.Printf("iio: size is less than or equal to 0")
		return
	}
	C.iio_buffer_set_data(b.handle, ptr)
}

// PushPartial will push the data written to the Buffer (from start to end) to the
func (b *Buffer) PushPartial(length int) (int, error) {
	i := C.iio_buffer_push_partial(b.handle, C.size_t(length))
	if i < 0 {
		return 0, syscall.Errno(-i)
	}
	return int(i), nil
}

// Push will push the data written to the Buffer (from start to end) to the device
func (b *Buffer) Push() (int, error) {
	i := C.iio_buffer_push(b.handle)
	if i < 0 {
		return 0, syscall.Errno(-i)
	}
	return int(i), nil
}

// Refill will fill the Buffer up with samples from the backing device.
func (b *Buffer) Refill() (int, error) {
	i := C.iio_buffer_refill(b.handle)
	if i < 0 {
		return 0, syscall.Errno(-i)
	}
	return int(i), nil
}

// Get Size gets the size of the buffer
func (b *Buffer) GetSize() int {
	return b.size
}

// Get Data returns the byte[] reference to the buffer
func (b *Buffer) GetData() []byte {
	return b.data
}

// Create Buffer creates a buffer and a buffer handle
func (d *Device) createBuffer(bytesCount int, cyclic bool) (*Buffer, error) {
	buf, err := C.iio_device_create_buffer(
		d.handle,
		C.size_t(bytesCount),
		C.bool(cyclic),
	)
	if buf == nil {
		return nil, err
	}
	var closed bool
	return &Buffer{
		handle: buf,
		closed: &closed,
		size:   bytesCount,
	}, nil
}

// CreateBuffer will create an iio_buffer from a given device.
// At least one channel must be Enabled prior to this call.
func (d *Device) CreateBuffer(bytesCount int) (*Buffer, error) {
	return d.createBuffer(bytesCount, false)
}

// CreateCyclicBuffer will create an iio_buffer from a given device.
// At least one channel must be Enabled prior to this call.
func (d *Device) CreateCyclicBuffer(bytesCount int) (*Buffer, error) {
	return d.createBuffer(bytesCount, true)
}

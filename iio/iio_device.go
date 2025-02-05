package iio

// #cgo pkg-config: libiio
// #include <iio.h>
// #include <stdlib.h>
import "C"
import (
	"unsafe"
)

// Device represents an IIO device
type Device struct {
	handle *C.struct_iio_device
}

// GetName returns the name of the device
func (dev *Device) GetName() string {
	name := C.iio_device_get_name(dev.handle)
	if name == nil {
		return ""
	}
	return C.GoString(name)
}

// GetID returns the ID of the device
func (dev *Device) GetID() string {
	return C.GoString(C.iio_device_get_id(dev.handle))
}

// GetChannelsCount returns the number of channels the device has
func (dev *Device) GetChannelsCount() uint {
	return uint(C.iio_device_get_channels_count(dev.handle))
}

// GetChannel returns the channel at the specified index
func (dev *Device) GetChannel(index uint) (*Channel, error) {
	handle := C.iio_device_get_channel(dev.handle, C.uint(index))
	if handle == nil {
		return nil, getLastError()
	}
	return &Channel{handle: handle}, nil
}

// GetAttr returns the value of the specified attribute
func (dev *Device) GetAttr(name string) (string, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var value *C.char
	ret := C.iio_device_attr_read(dev.handle, cName, &value)
	if ret < 0 {
		return "", getError(ret)
	}
	return C.GoString(value), nil
}

// SetAttr sets the value of the specified attribute
func (dev *Device) SetAttr(name, value string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	ret := C.iio_device_attr_write(dev.handle, cName, cValue)
	if ret < 0 {
		return getError(ret)
	}
	return nil
}

// GetDebugAttr returns the value of the specified debug attribute
func (dev *Device) GetDebugAttr(name string) (string, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var value *C.char
	ret := C.iio_device_debug_attr_read(dev.handle, cName, &value)
	if ret < 0 {
		return "", getError(ret)
	}
	return C.GoString(value), nil
}

// SetDebugAttr sets the value of the specified debug attribute
func (dev *Device) SetDebugAttr(name, value string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	ret := C.iio_device_debug_attr_write(dev.handle, cName, cValue)
	if ret < 0 {
		return getError(ret)
	}
	return nil
}

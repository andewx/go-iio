package iio

// #cgo pkg-config: libiio
// #include <iio.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

type DeviceStatus int

const (
	DeviceStatusUnknown DeviceStatus = iota
	DeviceStatusOK
	DeviceStatusError
	DeviceStatusConnected
	DeviceStatusDisconnected
)

// Device represents an IIO device
type Device struct {
	handle     *C.struct_iio_device
	channels   []*Channel
	attributes []string
	status     DeviceStatus
}

func NewDevice(handle *C.struct_iio_device) *Device {
	return &Device{handle: handle, status: DeviceStatusUnknown, channels: make([]*Channel, 0), attributes: make([]string, 0)}
}

func (dev *Device) Init() error {

	// Initialize channels
	dev.channels = make([]*Channel, dev.GetChannelsCount())
	for i := uint(0); i < dev.GetChannelsCount(); i++ {
		channel, err := dev.GetChannel(i)
		if err != nil {
			return err
		}
		dev.channels[i] = channel
	}

	// Initialize attributes
	dev.attributes = make([]string, dev.GetAttributesCount())
	for i := uint(0); i < dev.GetAttributesCount(); i++ {
		dev.attributes[i] = dev.GetAttr(i)
	}
	return nil
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

// GetAttributes returns the attributes of the device
func (dev *Device) GetAttributes() []string {
	return dev.attributes
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
	return NewChannel(handle), nil
}

// GetAttr returns the value of the specified attribute
func (dev *Device) GetAttr(name string) (string, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var value *C.char
	value = (*C.char)(C.malloc(C.size_t(512)))
	defer C.free(unsafe.Pointer(value))
	ret := C.iio_device_attr_read(dev.handle, cName, value, 512)
	if ret < 0 {
		return "", fmt.Errorf("iio: attribute %s not found", name)
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
		return fmt.Errorf("iio: attribute %s not found", name)
	}
	return nil
}

// GetDebugAttr returns the value of the specified debug attribute
func (dev *Device) GetDebugAttr(name string) (string, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var value *C.char
	value = (*C.char)(C.malloc(C.size_t(1024)))
	defer C.free(unsafe.Pointer(value))
	ret := C.iio_device_debug_attr_read(dev.handle, cName, value, 1024)
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
		return fmt.Errorf("iio: attribute %s not found", name)
	}
	return nil
}

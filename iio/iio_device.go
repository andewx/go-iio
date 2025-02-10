package iio

// #cgo pkg-config: libiio
// #include <iio.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/andewx/go-iio/common"
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
	attributes []*DeviceAttribute
	status     DeviceStatus
	_name      string
}

// NewDevice creates a new Device
func NewDevice(handle *C.struct_iio_device) *Device {
	dev := &Device{handle: handle, status: DeviceStatusUnknown, channels: nil, attributes: nil}
	dev._name = dev.GetName()
	dev.init()

	return dev
}

// Init - initialzes the device including reading in all channels and attributes.
func (dev *Device) init() error {
	common.PrintDebug(fmt.Sprintf("file iio_device.go::init--initiating|%s", dev._name), dev)
	dev.channels = GetChannels(dev)
	dev.attributes = GetDeviceAttributes(dev.handle)
	return nil
}

// Destroy destroys device handle and any channels and attributes
func (dev *Device) Destroy() {
	//Destroy any buffers and blocks
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
func (dev *Device) GetAttributes() []*DeviceAttribute {
	return dev.attributes
}

// GetChannelsCount returns the number of channels the device has
func (dev *Device) GetChannelsCount() uint {
	return uint(C.iio_device_get_channels_count(dev.handle))
}

// GetChannels returns the device channels
func (dev *Device) GetChannels() []*Channel {
	return dev.channels
}

// GetChannel returns the channel at the specified index
func (dev *Device) GetChannel(index uint) (*Channel, error) {
	handle := C.iio_device_get_channel(dev.handle, C.uint(index))
	if handle == nil {
		return nil, getLastError()
	}
	return NewChannel(handle, dev), nil
}

// GetAttr returns the value of the specified attribute
func (dev *Device) GetAttr(index int) (*DeviceAttribute, error) {
	if dev.attributes != nil {
		if index > len(dev.attributes) {
			return nil, fmt.Errorf("Error attribute index %d does not exist", index)
		}
		return dev.attributes[index], nil

	}
	return nil, fmt.Errorf("Device attribuets not initialized")
}

// SetAttr sets the value of a specified attribute for the device.
func (dev *Device) SetAttr(attrName string, value string) error {
	cAttrName := C.CString(attrName)
	defer C.free(unsafe.Pointer(cAttrName))

	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	// Assuming there's a C function to set device attributes
	if C.iio_device_attr_write(dev.handle, cAttrName, cValue) < 0 {
		return fmt.Errorf("failed to set attribute %s to %s", attrName, value)
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

// GetAttributesCount returns the number of attributes the device has
func (dev *Device) GetAttributesCount() int {
	return GetDeviceAttributesCount(dev.handle)
}

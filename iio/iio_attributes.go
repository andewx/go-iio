package iio

// #cgo pkg-config: libiio
// #include <iio.h>
// #include <stdlib.h>
import "C"

type ChannelAttributeType int

const (
	ChannelAttributeTypeString ChannelAttributeType = iota
	ChannelAttributeTypeInt
	ChannelAttributeTypeFloat
	ChannelAttributeTypeBool
	ChannelAttributeTypeLongLong
	ChannelAttributeTypeDouble
	ChannelAttributeTypeUnknown
)

// Channel Attributes are byte based attributes
type ChannelAttribute struct {
	data []byte
	typ  ChannelAttributeType
}

// Device Attributes are string based attributes

type DeviceAttribute struct {
	name string
}

func NewDeviceAttribute(name string) *DeviceAttribute {
	return &DeviceAttribute{name: name}
}

func (attr *DeviceAttribute) GetName() string {
	return attr.name
}

func GetAttributesCount(handle *C.struct_iio_device) int {
	return int(C.iio_device_get_attributes_count(handle))
}

func GetDeviceAttributes(handle *C.struct_iio_device) []*DeviceAttribute {
	count := GetAttributesCount(handle)
	attributes := make([]*Attribute, count)
	for i := 0; i < count; i++ {
		attributes[i] = NewDeviceAttribute(C.GoString(C.iio_device_get_attr(handle, C.uint(i))))
	}
	return attributes
}

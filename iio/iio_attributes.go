package iio

// #cgo pkg-config: libiio
// #include <iio.h>
// #include <stdlib.h>
import "C"
import "unsafe"

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

const (
	ATTR_DATA_DEFAULT_SIZE = 1024
)

// Channel Attributes are byte based attributes
type ChannelAttribute struct {
	name    string
	data    []byte
	ch_type ChannelAttributeType
}

func newChannelAttribute(name string, ch_type ChannelAttributeType) *ChannelAttribute {
	return &ChannelAttribute{name: name, ch_type: ch_type, data: make([]byte, ATTR_DATA_DEFAULT_SIZE)}
}

func GetChannelAttributesCount(handle *C.struct_iio_channel) int {
	return int(C.iio_channel_get_attrs_count(handle))
}

func GetChannelAttributes(handle *C.struct_iio_channel) []*ChannelAttribute {
	count := GetChannelAttributesCount(handle)
	attributes := make([]*ChannelAttribute, count)
	for i := 0; i < count; i++ {
		attributes[i] = newChannelAttribute(C.GoString(C.iio_channel_get_attr(handle, C.uint(i))), ChannelAttributeTypeString)
	}
	return attributes
}

func GetChannelAttributeByName(chAttributes []*ChannelAttribute, name string) *ChannelAttribute {
	for _, attr := range chAttributes {
		if attr.name == name {
			return attr
		}
	}
	return nil
}

func (attr *ChannelAttribute) GetChannelAttributeFilename(handle *C.struct_iio_channel, name string) string {
	return C.GoString(C.iio_channel_attr_get_filename(handle, C.CString(name)))
}

func (attr *ChannelAttribute) GetAttributeData(handle *C.struct_iio_channel, name string) []byte {
	var buffer *C.char
	buffer = (*C.char)(C.malloc(C.size_t(ATTR_DATA_DEFAULT_SIZE)))
	defer C.free(unsafe.Pointer(buffer))
	C.iio_channel_attr_read(handle, C.CString(name), buffer, ATTR_DATA_DEFAULT_SIZE)
	attr.data = C.GoBytes(unsafe.Pointer(buffer), ATTR_DATA_DEFAULT_SIZE)
	return attr.data
}

func (attr *ChannelAttribute) GetData() []byte {
	return attr.data
}

//-----------------------------------------------
// Device Attributes -- Device Attributes API consist of strings only
//-----------------------------------------------

type DeviceAttribute struct {
	name string
}

func newDeviceAttribute(name string) *DeviceAttribute {
	return &DeviceAttribute{name: name}
}

func (attr *DeviceAttribute) GetName() string {
	return attr.name
}

func GetDeviceAttributesCount(handle *C.struct_iio_device) int {
	return int(C.iio_device_get_attrs_count(handle))
}

func GetDeviceAttributes(handle *C.struct_iio_device) []*DeviceAttribute {
	count := GetDeviceAttributesCount(handle)
	attributes := make([]*DeviceAttribute, count)
	for i := 0; i < count; i++ {
		attributes[i] = newDeviceAttribute(C.GoString(C.iio_device_get_attr(handle, C.uint(i))))
	}
	return attributes
}

//-----------------------------------------------
// Context Attributes -- context attributes have both a name and value
//-----------------------------------------------

type ContextAttribute struct {
	name  string
	value string
}

func newContextAttribute(name, value string) *ContextAttribute {
	return &ContextAttribute{name: name, value: value}
}

func (attr *ContextAttribute) GetName() string {
	return attr.name
}

func (attr *ContextAttribute) GetValue() string {
	return attr.value
}

func GetContextAttributesCount(handle *C.struct_iio_context) int {
	return int(C.iio_context_get_attrs_count(handle))
}

func GetContextAttributes(handle *C.struct_iio_context) []*ContextAttribute {
	count := GetContextAttributesCount(handle)
	attributes := make([]*ContextAttribute, count)
	for i := 0; i < count; i++ {
		var name, value *C.char
		C.iio_context_get_attr(handle, C.uint(i), &name, &value)
		attributes[i] = newContextAttribute(C.GoString(name), C.GoString(value))
	}
	return attributes
}

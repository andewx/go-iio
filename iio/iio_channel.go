package iio

// #cgo pkg-config: libiio
// #include <iio.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

type ChannelStatus int

const (
	ChannelStatusUnknown ChannelStatus = iota
	ChannelStatusOK
	ChannelStatusError
)

// Channel represents an IIO channel
type Channel struct {
	handle     *C.struct_iio_channel
	attributes []*ChannelAttribute
	status     ChannelStatus
}

func NewChannel(handle *C.struct_iio_channel) *Channel {
	return &Channel{handle: handle, status: ChannelStatusUnknown, attributes: nil}
}

// Init() - After a handle for a channel is recieved we initiate the channel which merely reads in its attributes
func (ch *Channel) Init() error {
	var err error
	ch.attributes = GetChannelAttributes(ch.handle)
	if ch.attributes == nil {
		err = fmt.Errorf("Unable to establish IIO Connection to Channel and read attributes")
		return err
	}
	return nil
}

// GetID returns the ID of the channel
func (ch *Channel) GetID() string {
	return C.GoString(C.iio_channel_get_id(ch.handle))
}

// GetName returns the name of the channel
func (ch *Channel) GetName() string {
	name := C.iio_channel_get_name(ch.handle)
	if name == nil {
		return ""
	}
	return C.GoString(name)
}

// IsOutput returns true if the channel is an output channel
func (ch *Channel) IsOutput() bool {
	return bool(C.iio_channel_is_output(ch.handle))
}

// IsScanned returns true if the channel is scanned
func (ch *Channel) IsScanned() bool {
	return bool(C.iio_channel_is_scan_element(ch.handle))
}

// Enable enables the channel
func (ch *Channel) Enable() {
	C.iio_channel_enable(ch.handle)
}

// Disable disables the channel
func (ch *Channel) Disable() {
	C.iio_channel_disable(ch.handle)
}

// IsEnabled returns true if the channel is enabled
func (ch *Channel) IsEnabled() bool {
	return bool(C.iio_channel_is_enabled(ch.handle))
}

// GetAttr returns the value of the specified attribute
func (ch *Channel) GetAttr(name string) (string, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var value *C.char
	value = (*C.char)(C.malloc(C.size_t(256)))
	defer C.free(unsafe.Pointer(value))
	ret := C.iio_channel_attr_read(ch.handle, cName, value, 256)
	if ret < 0 {
		return "", getError(ret)
	}
	return C.GoString(value), nil
}

// SetAttr sets the value of the specified attribute read the iio.h specification
// For the multi-attribute write method using 32 bit signed integer headers for block data
func (ch *Channel) SetAttr(name, value string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	ret := C.iio_channel_attr_write(ch.handle, cName, cValue)
	if ret < 0 {
		return getError(ret)
	}
	return nil
}

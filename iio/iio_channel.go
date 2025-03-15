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

type ChannelStatus int

const (
	// Unknown status
	ChannelStatusUnknown  ChannelStatus = iota // Unknown status
	ChannelStatusOK                            // OK status
	ChannelStatusError                         // Error status
	ChannelStatusDisabled                      // Disabled status
	ChannelStatusEnabled                       // Enabled status
)

const (
	KB = 1024      // 1024 bytes
	MB = 1024 * KB // 1024 * 1024 bytes
)

type Endianess int

const (
	EndianessLittle Endianess = iota // Little endian
	EndianessBig                     // Big endian
)

// DataFormat represents the format of the data in the buffer
type DataFormat struct {
	bits  int
	block int
	end   Endianess
}

// Channel represents an IIO channel
type Channel struct {
	handle     *C.struct_iio_channel
	attributes []*ChannelAttribute
	status     ChannelStatus
	device     *Device
	buffer     *Buffer
	output     bool
	_name      string
	_id        string
	format     DataFormat
}

// NewChannel - Creates a new channel
func NewChannel(handle *C.struct_iio_channel, dev *Device) *Channel {
	ch := &Channel{handle: handle, status: ChannelStatusUnknown, attributes: nil, device: dev, format: DataFormat{bits: 12, block: 16, end: EndianessLittle}}
	ch.init()
	return ch
}

// Init - After a handle for a channel is recieved we initiate the channel which merely reads in its attributes
func (ch *Channel) init() error {
	ch.attributes = GetChannelAttributes(ch.handle)
	ch.output = ch.IsOutput()
	ch._name = ch.name()
	ch._id = ch.id()
	ch.status = ChannelStatusOK
	common.PrintDebug(fmt.Sprintf("file iio_channel.go::init|%s", ch._name), ch)
	return nil
}

// initBuffer initializes a buffer for the channel
func (ch *Channel) initBuffer(sizeKb int) error {
	ch.Enable()
	buf, err := ch.device.createBuffer(sizeKb*KB, false)
	if err != nil {
		return err
	}
	ch.buffer = buf
	return nil
}

// Open opens a buffer for the channel
func (ch *Channel) Open(size int) error {
	err := ch.initBuffer(size)
	return err
}

// GetID returns the ID of the channel
func (ch *Channel) id() string {
	str := C.GoString(C.iio_channel_get_id(ch.handle))
	return str
}

func (ch *Channel) name() string {
	str := C.GoString(C.iio_channel_get_name(ch.handle))
	if str == "" {
		return ch.id()
	}
	return str
}

// DirectionalName returns the name of the channel with the direction prefix
func (ch *Channel) DirectionalName() string {
	if ch.output {
		return "out_" + ch.Name()
	}
	return "in_" + ch.Name()
}

// Name returns the name of the channel
func (ch *Channel) Name() string {
	return ch._name
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
	ch.status = ChannelStatusEnabled
}

// Disable disables the channel
func (ch *Channel) Disable() {
	C.iio_channel_disable(ch.handle)
	ch.status = ChannelStatusDisabled
}

// Close closes the channel and removes the buffer
func (ch *Channel) Close() {
	ch.Disable()
	ch.buffer.Close()
}

func (ch *Channel) Write(data []byte) {
	ch.buffer.Write(unsafe.Pointer(&data[0]), len(data))
}

func (ch *Channel) Read(data []byte) {
	ch.buffer.Read(ch, unsafe.Pointer(&data[0]), len(data))
}

// IsEnabled returns true if the channel is enabled
func (ch *Channel) IsEnabled() bool {
	return bool(C.iio_channel_is_enabled(ch.handle))
}

// GetChannelsCount returns the number of channels in the device
func GetChannelsCount(dev *Device) int {
	return int(C.iio_device_get_channels_count(dev.handle))
}

// GetChannels returns all channels in the device
func GetChannels(dev *Device) map[string]*Channel {
	count := GetChannelsCount(dev)
	channels := make([]*Channel, count)
	channelsMap := make(map[string]*Channel)
	for i := 0; i < count; i++ {
		channels[i] = NewChannel(C.iio_device_get_channel(dev.handle, C.uint(i)), dev)
		channels[i].init()
		channelsMap[channels[i]._id] = channels[i]
	}
	return channelsMap
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

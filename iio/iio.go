package iio

// #cgo pkg-config: libiio
// #include <iio.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

// DeviceConfig represents high-level configuration for an IIO device
type DeviceConfig struct {
	SampleRate  int    // Hz
	BufferSize  uint   // samples
	Channels    []bool // enabled channels
	IsCyclic    bool   // cyclic buffer mode
	TriggerName string // optional hardware trigger
}

// Stream represents a continuous data stream from/to an IIO device
type Stream struct {
	device  *Device
	buffer  *Buffer
	samples []byte
	config  DeviceConfig
}

// func Verision - Returns version of the LibIIO library this was compiled against.
func Version() string {
	return "v-0.26"
}

// NewLocalDevice creates and configures a local IIO device
func NewLocalDevice(deviceName string, config DeviceConfig) (*Context, *Device, error) {
	ctx, err := CreateContext()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create local context: %w", err)
	}

	dev, err := setupDevice(ctx, deviceName, config)
	return ctx, dev, err
}

// NewNetworkDevice creates and configures a remote IIO device
func NewNetworkDevice(hostname, deviceName string, config DeviceConfig) (*Context, *Device, error) {
	ctx, err := CreateNetworkContext(hostname)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create network context: %w", err)
	}
	dev, err := setupDevice(ctx, deviceName, config)
	return ctx, dev, err
}

// NewXMLDevice creates and configures an IIO device from XML description
func NewXMLDevice(xmlPath, deviceName string, config DeviceConfig) (*Context, *Device, error) {
	ctx, err := CreateXMLContext(xmlPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create XML context: %w", err)
	}
	dev, err := setupDevice(ctx, deviceName, config)
	return ctx, dev, err
}

// setupDevice configures an IIO device according to the provided config
func setupDevice(ctx *Context, deviceName string, config DeviceConfig) (*Device, error) {
	dev, err := ctx.FindDevice(deviceName)
	if err != nil {
		ctx.Close()
		return nil, fmt.Errorf("%s device not found: %w", deviceName, err)
	}

	err = dev.Init()
	if err != nil {
		return nil, fmt.Errorf("%s device couldn't be initiated. Returned with error %w", deviceName, err)
	}

	// Configure channels
	for i := 0; i < len(dev.channels); i++ {
		c := dev.channels[i]
		c.Init()
		if err != nil {
			return nil, fmt.Errorf("%s device channel not initiated %w", err)
		}
	}

	// Create buffer
	buf, err := dev.CreateBuffer(int(config.BufferSize))
	if err != nil {
		ctx.Close()
		return nil, fmt.Errorf("failed to create buffer: %w", err)
	}

	strm := &Stream{
		device:  dev,
		buffer:  buf,
		config:  config,
		samples: make([]byte, buf.GetSize()),
	}

	dev.stream = strm
	return dev, nil
}

// Read reads samples from the device into the provided slice
// The slice must be of a supported type (int16, int32, float32, float64)
func (s *Stream) Read(samples interface{}) (int, error) {
	// Validate input slice
	val := reflect.ValueOf(samples)
	if val.Kind() != reflect.Slice {
		return 0, fmt.Errorf("samples must be a slice")
	}

	// Refill the buffer
	_, err := s.buffer.Refill()
	if err != nil {
		return 0, fmt.Errorf("failed to refill buffer: %w", err)
	}

	// Get buffer data
	data := s.buffer.GetData()
	if len(data) == 0 {
		return 0, nil
	}

	// Calculate number of samples to copy
	sampleSize := val.Type().Elem().Size()
	numSamples := len(data) / int(sampleSize)
	if numSamples > val.Len() {
		numSamples = val.Len()
	}

	// Copy data using reflection
	header := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&data[0])),
		Len:  numSamples,
		Cap:  numSamples,
	}

	src := reflect.NewAt(val.Type(), unsafe.Pointer(&header)).Elem()
	reflect.Copy(val, src)

	return numSamples, nil
}

// Write writes samples from the provided slice to the device
// The slice must be of a supported type (int16, int32, float32, float64)
func (s *Stream) Write(samples interface{}) (int, error) {
	// Validate input slice
	val := reflect.ValueOf(samples)
	if val.Kind() != reflect.Slice {
		return 0, fmt.Errorf("samples must be a slice")
	}

	// Get buffer data
	data := s.buffer.GetData()
	if len(data) == 0 {
		return 0, nil
	}

	// Calculate number of samples to copy
	sampleSize := val.Type().Elem().Size()
	numSamples := len(data) / int(sampleSize)
	if numSamples > val.Len() {
		numSamples = val.Len()
	}

	// Copy data using reflection
	header := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&data[0])),
		Len:  numSamples,
		Cap:  numSamples,
	}

	dst := reflect.NewAt(val.Type(), unsafe.Pointer(&header)).Elem()
	reflect.Copy(dst, val)

	// Push the buffer
	_, err := s.buffer.Push()
	if err != nil {
		return 0, fmt.Errorf("failed to push buffer: %w", err)
	}

	return numSamples, nil
}

// Close releases all resources associated with the stream
func (s *Stream) Close() error {
	if s.buffer != nil {
		s.buffer.Close()
		s.buffer = nil
	}
	return nil
}

// Example usage:
/*
func main() {
    // Configure device
    config := DeviceConfig{
        SampleRate:  1000,
        BufferSize:  1024,
        Channels:    []bool{true, true, false, false},
        IsCyclic:    false,
        TriggerName: "trigger0",
    }

    // Create local device
    stream, err := NewLocalDevice("device0", config)
    if err != nil {
        log.Fatal(err)
    }
    defer stream.Close()

    // Read samples
    samples := make([]int16, 1024)
    n, err := stream.Read(samples)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Read %d samples\n", n)

    // Write samples
    n, err = stream.Write(samples)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Wrote %d samples\n", n)
}
*/

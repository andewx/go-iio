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

// Stream represents a continuous data stream from/to an IIO device
type Stream struct {
	device  *Device
	buffer  *Buffer
	samples []byte
	config  DeviceConfig
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

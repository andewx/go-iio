package iio

// #cgo pkg-config: libiio
// #include <iio.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"

	"github.com/andewx/go-iio/common"
)

// DeviceConfig represents high-level configuration for an IIO device
type DeviceConfig struct {
	SampleRate  int    // Hz
	BufferSize  uint   // samples
	Channels    []bool // enabled channels
	IsCyclic    bool   // cyclic buffer mode
	TriggerName string // optional hardware trigger
}

// func Verision - Returns version of the LibIIO library this was compiled against.
func Version() string {
	return "v-0.26"
}

// NewLocalDevice creates and configures a local IIO device
func NewLocalDevice(deviceName string) (*Context, error) {
	ctx, err := CreateContext()
	if err != nil {
		return nil, fmt.Errorf("failed to create local context: %w", err)
	}
	return ctx, err
}

// NewNetworkDevice creates and configures a remote IIO device
func NewNetworkDevice(hostname string) (*Context, error) {
	common.PrintDebug("File:iio.go // Line:39 // Function:NewNetworkDevice // Creating network context", hostname)
	ctx, err := CreateNetworkContext(hostname)
	if err != nil {
		return nil, fmt.Errorf("failed to create network context: %w", err)
	}
	return ctx, err
}

// NewXMLDevice creates and configures an IIO device from XML description
func NewXMLDevice(xmlPath string) (*Context, error) {
	ctx, err := CreateXMLContext(xmlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create XML context: %w", err)
	}
	return ctx, err
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

package sdr

import (
	"fmt"
	"math"

	"github.com/andewx/go-iio/common"
	"github.com/andewx/go-iio/iio"
)

// ModulationType defines the type of modulation
type ModulationType int

// Modulation
const (
	ModulationAM  ModulationType = iota //ModulationAM
	ModulationFM                        //ModulationFM
	ModulationPSK                       //ModluationPSK
	ModulationQAM                       //Modulation QAM
)

// PHYDevice - Is a Physcial Layer Device implementation of an SDR Radio
type PHYDevice interface {
	Connect() error
	Close() error
	String() string
}

// Config represents SDR-specific configuration
type Config struct {
	Modulation       ModulationType
	SymbolFrequency  float64 // Hz
	Alpha            float64 // Root-raised cosine filter alpha
	CarrierFrequency float64 // Hz
	TXFrequency      float64 // Hz
	RXFrequency      float64 // Hz
	RFBandwidth      float64 // Hz
	SampleRate       float64 // Calculated Baseband Sampling Rate
	FMCoeff          float64 // Modulation index for FM
	QAMOrder         int     // QAM constellation order (4, 16, 64, etc.)
	PSKOrder         int     // PSK constellation order (2, 4, 8, etc.
	BitDepth         int     // Bits Per Sample
	FourierPoint     int     //FFT Size
	Gain             float64 // Db
}

/*
SDR struct
ctx: iio context
stream: iio stream
config: Config
name: device name

Represents a software defined radio device and especially its context and high level configuration.
*/
type SDR struct {
	ctx      *iio.Context
	config   Config
	name     string
	Platform PHYDevice
}

// ConnectNetwork - connects to a specified host context and sets up an SDR parameters bench
func ConnectNetwork(host string, name string) (*SDR, error) {
	common.PrintDebug("File:sdr.go // Line:66 // Function:ConnectNetwork // Connecting to Network Device", host)
	var s SDR
	s = SDR{config: DefaultConfig(), name: name}
	ctx, err := iio.NewNetworkDevice(host) //q: how does NewNetworkDevice setup a stream..answer it doesn't we just get a stream high level wrapper struct
	if err != nil {
		return nil, err
	}
	s.ctx = ctx

	var switchErr error
	switch name {
	case "ad9361":
		s.Platform, switchErr = newAD9361(ctx)
		if switchErr != nil {
			return nil, switchErr
		}
	default:
		return nil, fmt.Errorf("invalid device name: %s", name)
	}

	return &s, nil
}

// DefaultConfig returns a default Config with typical settings.
func DefaultConfig() Config {
	return Config{
		CarrierFrequency: 100e6, // Default frequency: 100 MHz
		TXFrequency:      500e3, // Default Transmit Frequency 500 kHz
		RXFrequency:      500e3, // Default Receive Frequency 500 kHz
		SampleRate:       1e6,   // Default sample rate: 1 MS/s
		Gain:             10.0,  // Default gain: 10 dB
		Alpha:            0.5,   // Default RRC Filter Alpha Parameter 0.5
		RFBandwidth:      1e6,   // Default RF Bandwidth 1MhZ
		FMCoeff:          0.7,   // FM Modulation index coefficient
		QAMOrder:         4,     // QAM Default
		PSKOrder:         2,     // PSK Default
		FourierPoint:     1024,  // FFT Size Default 1024
		SymbolFrequency:  1e6,   // Default Symbol Frequency 1 MhZ.
	}
}

// GenerateRRCFilter generates Root-Raised Cosine filter taps
func GenerateRRCFilter(numTaps int, alpha float64, symbolsPerTap float64) []float64 {
	taps := make([]float64, numTaps)

	for i := 0; i < numTaps; i++ {
		t := float64(i-numTaps/2) / symbolsPerTap
		if t == 0.0 {
			taps[i] = 1.0 + alpha*(4.0/math.Pi-1.0)
		} else if math.Abs(t) == 1.0/(4.0*alpha) {
			taps[i] = alpha / math.Sqrt(2.0) * ((1.0+2.0/math.Pi)*
				math.Sin(math.Pi/(4.0*alpha)) + (1.0-2.0/math.Pi)*
				math.Cos(math.Pi/(4.0*alpha)))
		} else {
			taps[i] = (math.Sin(math.Pi*t*(1.0-alpha)) +
				4.0*alpha*t*math.Cos(math.Pi*t*(1.0+alpha))) /
				(math.Pi * t * (1.0 - (4.0*alpha*t)*(4.0*alpha*t)))
		}
	}

	return taps
}

// Print prints the SDR device tree along with all device attributes, channels, and channel attributes
func (sdr *SDR) Print() {
	common.PrintDebug("File:sdr.go // Line:119 // Function:Print // Printing SDR Device Tree", sdr)
	fmt.Printf("SDR Platform %s\n%s", sdr.name, sdr.Platform.String())
}

package devices

import (
	"math"

	"github.com/andewx/go-iio/iio"
)

// ModulationType defines the type of modulation
type ModulationType int

const (
	ModulationAM ModulationType = iota
	ModulationFM
	ModulationPSK
	ModulationQAM
)

// SDRConfig represents SDR-specific configuration
type SDRConfig struct {
	Modulation  ModulationType
	SymbolRate  float64
	FilterAlpha float64 // Root-raised cosine filter alpha
	FrequencyOffset float64 // Baseband - Carrier frequency offset
	CarrierFrequency float64 //Carrier Frequency Target Zero-IF
	SampleRate float64  // Calculated Baseband Sampling Rate
	ModIndex    float64 // Modulation index for FM
	QAMOrder    int     // QAM constellation order (4, 16, 64, etc.)
	PSKOrder    int     // PSK constellation order (2, 4, 8, etc.)
}

type SDR struct {
	ctx *iio.Context
	dev *iio.Device
	rx  *iio.Device
	tx  *iio.Device
	buf *iio.Buffer
	config SDRConfig
	name string
}


func NewNetworkSDR(host string. deviceName string) (*SDR, error) {
	ctx, err := iio.CreateNetworkContext(host)
	if err != nil {
		return nil, err
	}
	return &SDR{ctx: ctx, name: deviceName}, nil
}

func (s *SDR)ConfigureSDR(symbolRate float64, nyquistRate float64, alphaFreq float64, carrierFreq float64, modIndex float64, qamOrder uint, pskOrder uint ){
	if nyquistRate < 2.0{
		nyquistRate = 2.0
	}

	if modIndex <= 0 || modIndex >= 1.0{
		modIndex = 0.5	
	}

	if carrierFrequency <= 0 {
		carrierFrequency = 1000
	}

	if alphaFreq <= 0 || alphaFreq >= 1.0{
		alphaFreq = 0.5
	}
	config := SDRConfig{SymbolRate:symbolRate, FilterAlpha:alphaFreq, CarrierFrequency:carrierFreq, SampleRate: nyquistRate*(symbolRate + FrequencyOffset), ModIndex: modIndex; QAMOrder:qamOrder; PSKOrder:pskOrder }
	s.config = config
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

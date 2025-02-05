package devices

import (
	"fmt"
	"math"
	"math/cmplx"
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
	CarrierFreq float64 // Carrier frequency offset
	ModIndex    float64 // Modulation index for FM
	QAMOrder    int     // QAM constellation order (4, 16, 64, etc.)
	PSKOrder    int     // PSK constellation order (2, 4, 8, etc.)
}

// ConfigureDDS sets up the Direct Digital Synthesis
func (dev *AD9361Device) ConfigureDDS() error {
	if err := dev.tx.SetAttr("out_altvoltage0_TX1_I_F1_frequency",
		fmt.Sprintf("%d", dev.config.DDSFrequency)); err != nil {
		return err
	}

	scale := fmt.Sprintf("%.6f", dev.config.DDSScale)
	if err := dev.tx.SetAttr("out_altvoltage0_TX1_I_F1_scale", scale); err != nil {
		return err
	}

	return nil
}

// SetFIRFilter configures custom FIR filter taps
func (dev *AD9361Device) SetFIRFilter(taps []int16) error {
	// Convert taps to string
	tapStr := ""
	for i, tap := range taps {
		if i > 0 {
			tapStr += ","
		}
		tapStr += fmt.Sprintf("%d", tap)
	}

	return dev.phy.SetAttr("in_voltage_filter_fir_en", tapStr)
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

// Modulate applies the specified modulation to the input samples
func (dev *AD9361Device) Modulate(samples []complex128, config SDRConfig) ([]complex128, error) {
	output := make([]complex128, len(samples))

	switch config.Modulation {
	case ModulationAM:
		return dev.modulateAM(samples, config)
	case ModulationFM:
		return dev.modulateFM(samples, config)
	case ModulationPSK:
		return dev.modulatePSK(samples, config)
	case ModulationQAM:
		return dev.modulateQAM(samples, config)
	default:
		return nil, fmt.Errorf("unsupported modulation type")
	}
}

func (dev *AD9361Device) modulateAM(samples []complex128, config SDRConfig) ([]complex128, error) {
	output := make([]complex128, len(samples))
	carrier := complex(math.Cos(config.CarrierFreq), math.Sin(config.CarrierFreq))

	for i, sample := range samples {
		// AM modulation: (1 + m*signal)*carrier
		magnitude := 1.0 + real(sample)
		output[i] = complex(magnitude*real(carrier), magnitude*imag(carrier))
	}

	return output, nil
}

func (dev *AD9361Device) modulateFM(samples []complex128, config SDRConfig) ([]complex128, error) {
	output := make([]complex128, len(samples))
	phase := 0.0

	for i, sample := range samples {
		// FM modulation: exp(j*(wc*t + m*integral(signal)))
		phase += config.ModIndex * real(sample)
		output[i] = cmplx.Rect(1.0, phase)
	}

	return output, nil
}

func (dev *AD9361Device) modulatePSK(samples []complex128, config SDRConfig) ([]complex128, error) {
	output := make([]complex128, len(samples))
	symbolAngle := 2 * math.Pi / float64(config.PSKOrder)

	for i, sample := range samples {
		// PSK modulation: map symbols to constellation points
		symbol := int(real(sample)) % config.PSKOrder
		angle := float64(symbol) * symbolAngle
		output[i] = cmplx.Rect(1.0, angle)
	}

	return output, nil
}

func (dev *AD9361Device) modulateQAM(samples []complex128, config SDRConfig) ([]complex128, error) {
	output := make([]complex128, len(samples))
	sqrtM := int(math.Sqrt(float64(config.QAMOrder)))
	scale := 2.0 * float64(sqrtM-1)

	for i, sample := range samples {
		// QAM modulation: map symbols to square constellation
		symbol := int(real(sample)) % config.QAMOrder
		I := float64(symbol%sqrtM)*2.0 - float64(sqrtM-1)
		Q := float64(symbol/sqrtM)*2.0 - float64(sqrtM-1)
		output[i] = complex(I/scale, Q/scale)
	}

	return output, nil
}

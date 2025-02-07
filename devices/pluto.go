package devices

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/andewx/go-iio/iio"
)

type SDRDevice interface {
	Connect(host string, name string, config SDRConfig, enableChannels []bool) error
	Close() error
	Modulate(samples []complex64, config SDRConfig) ([]complex64, error)
}

// AD9361Config represents the configuration for AD9361 transceiver
type AD9361 struct {

	// Base Configuration
	Context *SDR

	// RF Configuration
	RXFrequency uint64  // Hz
	TXFrequency uint64  // Hz
	RFBandwidth uint32  // Hz
	GainControl string  // "manual", "slow_attack", "fast_attack"
	RXGain      float64 // dB
	TXGain      float64 // dB

	// Advanced Configuration
	EnableDDS          bool    // Enable Direct Digital Synthesis
	DDSFrequency       uint32  // DDS frequency in Hz
	DDSScale           float64 // DDS scale factor
	EnableFIR          bool    // Enable FIR filtering
	FIRTaps            []int16 // Custom FIR filter taps
	QuadratureTracking bool    // Enable quadrature tracking
	DCOffsetTracking   bool    // Enable DC offset tracking
}

// NewAD9361 creates and configures an AD9361 device
func NewAD9361(hostname string, config AD9361Config) (*AD9361Device, error) {
	// Create network context
	ctx, err := iio.CreateNetworkContext(hostname)
	if err != nil {
		return nil, fmt.Errorf("failed to create context: %w", err)
	}

	// Find devices
	phy, err := ctx.FindDevice("ad9361-phy")
	if err != nil {
		return nil, fmt.Errorf("failed to find PHY device: %w", err)
	}

	rx, err := ctx.FindDevice("cf-ad9361-lpc")
	if err != nil {
		return nil, fmt.Errorf("failed to find RX device: %w", err)
	}

	tx, err := ctx.FindDevice("cf-ad9361-dds-core-lpc")
	if err != nil {
		return nil, fmt.Errorf("failed to find TX device: %w", err)
	}

	dev := &AD9361Device{
		phy:    phy,
		rx:     rx,
		tx:     tx,
		config: config,
	}

	// Configure the device
	if err := dev.Configure(); err != nil {
		return nil, err
	}

	return dev, nil
}

// Configure applies the configuration to the AD9361 device
func (dev *AD9361Device) Configure() error {
	// Configure RF parameters
	if err := dev.SetRXFrequency(dev.config.RXFrequency); err != nil {
		return err
	}
	if err := dev.SetTXFrequency(dev.config.TXFrequency); err != nil {
		return err
	}
	if err := dev.SetSampleRate(dev.config.SampleRate); err != nil {
		return err
	}
	if err := dev.SetRFBandwidth(dev.config.RFBandwidth); err != nil {
		return err
	}

	// Configure gain control
	if err := dev.SetGainControl(dev.config.GainControl); err != nil {
		return err
	}
	if err := dev.SetRXGain(dev.config.RXGain); err != nil {
		return err
	}
	if err := dev.SetTXGain(dev.config.TXGain); err != nil {
		return err
	}

	// Configure tracking options
	if err := dev.SetQuadratureTracking(dev.config.QuadratureTracking); err != nil {
		return err
	}
	if err := dev.SetDCOffsetTracking(dev.config.DCOffsetTracking); err != nil {
		return err
	}

	// Configure FIR filter if enabled
	if dev.config.EnableFIR {
		if err := dev.SetFIRFilter(dev.config.FIRTaps); err != nil {
			return err
		}
	}

	// Configure DDS if enabled
	if dev.config.EnableDDS {
		if err := dev.ConfigureDDS(); err != nil {
			return err
		}
	}

	return nil
}

// SetRXFrequency sets the RX frequency
func (dev *AD9361Device) SetRXFrequency(freq uint64) error {
	return dev.phy.SetAttr("out_altvoltage0_RX_LO_frequency", fmt.Sprintf("%d", freq))
}

// SetTXFrequency sets the TX frequency
func (dev *AD9361Device) SetTXFrequency(freq uint64) error {
	return dev.phy.SetAttr("out_altvoltage1_TX_LO_frequency", fmt.Sprintf("%d", freq))
}

// SetSampleRate sets the sample rate for both RX and TX
func (dev *AD9361Device) SetSampleRate(rate uint32) error {
	if err := dev.phy.SetAttr("in_voltage_sampling_frequency", fmt.Sprintf("%d", rate)); err != nil {
		return err
	}
	return dev.phy.SetAttr("out_voltage_sampling_frequency", fmt.Sprintf("%d", rate))
}

// SetRFBandwidth sets the RF bandwidth for both RX and TX
func (dev *AD9361Device) SetRFBandwidth(bw uint32) error {
	if err := dev.phy.SetAttr("in_voltage_rf_bandwidth", fmt.Sprintf("%d", bw)); err != nil {
		return err
	}
	return dev.phy.SetAttr("out_voltage_rf_bandwidth", fmt.Sprintf("%d", bw))
}

// SetGainControl sets the gain control mode
func (dev *AD9361Device) SetGainControl(mode string) error {
	return dev.phy.SetAttr("in_voltage0_gain_control_mode", mode)
}

// SetRXGain sets the RX gain in dB
func (dev *AD9361Device) SetRXGain(gain float64) error {
	return dev.phy.SetAttr("in_voltage0_hardwaregain", fmt.Sprintf("%.2f", gain))
}

// SetTXGain sets the TX gain in dB
func (dev *AD9361Device) SetTXGain(gain float64) error {
	return dev.phy.SetAttr("out_voltage0_hardwaregain", fmt.Sprintf("%.2f", gain))
}

// SetQuadratureTracking enables/disables quadrature tracking
func (dev *AD9361Device) SetQuadratureTracking(enable bool) error {
	val := "0"
	if enable {
		val = "1"
	}
	return dev.phy.SetAttr("in_voltage0_quadrature_tracking_en", val)
}

// SetDCOffsetTracking enables/disables DC offset tracking
func (dev *AD9361Device) SetDCOffsetTracking(enable bool) error {
	val := "0"
	if enable {
		val = "1"
	}
	return dev.phy.SetAttr("in_voltage0_dc_offset_tracking_en", val)
}

// Close releases all resources associated with the device
func (dev *AD9361Device) Close() error {
	if dev.stream != nil {
		dev.stream.Close()
	}
	return nil
}

// Modulate applies the specified modulation to the input samples
func (dev *AD9361Device) Modulate(samples []complex64, config SDRConfig) ([]complex64, error) {

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

func (dev *AD9361Device) modulateAM(samples []complex64, config SDRConfig) ([]complex64, error) {
	output := make([]complex64, len(samples))
	carrier := complex(math.Cos(config.CarrierFreq), math.Sin(config.CarrierFreq))

	for i, sample := range samples {
		// AM modulation: (1 + m*signal)*carrier
		magnitude := float32(1.0 + real(sample))
		output[i] = complex(magnitude*float32(real(carrier)), magnitude*float32(imag(carrier)))
	}

	return output, nil
}

func (dev *AD9361Device) modulateFM(samples []complex64, config SDRConfig) ([]complex64, error) {
	output := make([]complex64, len(samples))
	phase := 0.0

	for i, sample := range samples {
		// FM modulation: exp(j*(wc*t + m*integral(signal)))
		phase += config.ModIndex * float64(real(complex128(sample)))
		output[i] = complex64(cmplx.Rect(1.0, phase))
	}

	return output, nil
}

func (dev *AD9361Device) modulatePSK(samples []complex64, config SDRConfig) ([]complex64, error) {
	output := make([]complex64, len(samples))
	symbolAngle := 2 * math.Pi / float64(config.PSKOrder)

	for i, sample := range samples {
		// PSK modulation: map symbols to constellation points
		symbol := int(real(sample)) % config.PSKOrder
		angle := float64(symbol) * symbolAngle
		output[i] = complex64(cmplx.Rect(1.0, angle))
	}

	return output, nil
}

func (dev *AD9361Device) modulateQAM(samples []complex64, config SDRConfig) ([]complex64, error) {
	output := make([]complex64, len(samples))
	sqrtM := int(math.Sqrt(float64(config.QAMOrder)))
	scale := 2.0 * float64(sqrtM-1)

	for i, sample := range samples {
		// QAM modulation: map symbols to square constellation
		symbol := int(real(sample)) % config.QAMOrder
		I := float64(symbol%sqrtM)*2.0 - float64(sqrtM-1)
		Q := float64(symbol/sqrtM)*2.0 - float64(sqrtM-1)
		output[i] = complex64(complex(I/scale, Q/scale))
	}

	return output, nil
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

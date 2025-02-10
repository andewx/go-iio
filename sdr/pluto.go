package sdr

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/andewx/go-iio/common"
	"github.com/andewx/go-iio/iio"
)

// AD9361 represents AD9361 Device
type AD9361 struct {

	//Devices
	ctx *iio.Context
	phy *iio.Device
	rx  *iio.Device
	tx  *iio.Device
	sdr *SDR

	// RF Configuration
	RXFrequency      uint64  // Hz
	TXFrequency      uint64  // Hz
	RFBandwidth      uint32  // Hz
	CarrierFrequency float64 //Hz
	GainControl      string  // "manual", "slow_attack", "fast_attack"
	RXGain           float64 // dB
	TXGain           float64 // dB
	SampleRate       float64 // SPS

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
func newAD9361(ctx *iio.Context) (*AD9361, error) {

	common.PrintDebug("File:pluto.go // Line:73 // Function:newAD9361 // Creating AD9361 device", ctx)

	// Find devices
	phy, err := ctx.FindDevice("ad9361-phy")
	if err != nil {
		return nil, fmt.Errorf("failed to find PHY device: %w", err)
	}

	common.PrintValid("File:pluto.go // Line:54 // Function:newAD9361 // PHY device found")

	rx, err := ctx.FindDevice("cf-ad9361-lpc")
	if err != nil {
		return nil, fmt.Errorf("failed to find RX device: %w", err)
	}

	common.PrintValid("File:pluto.go // Line:54 // Function:newAD9361 // RX device found")

	tx, err := ctx.FindDevice("cf-ad9361-dds-core-lpc")
	if err != nil {
		return nil, fmt.Errorf("failed to find TX device: %w", err)
	}

	dev := &AD9361{
		ctx: ctx,
		phy: phy,
		rx:  rx,
		tx:  tx,
	}

	return dev, nil
}

// Configure applies the configuration to the AD9361 device
func (dev *AD9361) Configure() error {
	// Configure RF parameters
	if err := dev.SetRXFrequency(dev.RXFrequency); err != nil {
		return err
	}
	if err := dev.SetTXFrequency(dev.TXFrequency); err != nil {
		return err
	}
	if err := dev.SetSampleRate(dev.SampleRate); err != nil {
		return err
	}
	if err := dev.SetRFBandwidth(dev.RFBandwidth); err != nil {
		return err
	}

	// Configure gain control
	if err := dev.SetGainControl(dev.GainControl); err != nil {
		return err
	}
	if err := dev.SetRXGain(dev.RXGain); err != nil {
		return err
	}
	if err := dev.SetTXGain(dev.TXGain); err != nil {
		return err
	}

	// Configure tracking options
	if err := dev.SetQuadratureTracking(dev.QuadratureTracking); err != nil {
		return err
	}
	if err := dev.SetDCOffsetTracking(dev.DCOffsetTracking); err != nil {
		return err
	}

	// Configure FIR filter if enabled
	if dev.EnableFIR {
		if err := dev.SetFIRFilter(dev.FIRTaps); err != nil {
			return err
		}
	}

	// Configure DDS if enabled
	if dev.EnableDDS {
		if err := dev.ConfigureDDS(); err != nil {
			return err
		}
	}

	return nil
}

// SetRXFrequency sets the RX frequency
func (dev *AD9361) SetRXFrequency(freq uint64) error {
	return dev.phy.SetAttr("out_altvoltage0_RX_LO_frequency", fmt.Sprintf("%d", freq))
}

// SetTXFrequency sets the TX frequency
func (dev *AD9361) SetTXFrequency(freq uint64) error {
	return dev.phy.SetAttr("out_altvoltage1_TX_LO_frequency", fmt.Sprintf("%d", freq))
}

// SetSampleRate sets the sample rate for both RX and TX
func (dev *AD9361) SetSampleRate(rate float64) error {
	if err := dev.phy.SetAttr("in_voltage_sampling_frequency", fmt.Sprintf("%d", rate)); err != nil {
		return err
	}
	return dev.phy.SetAttr("out_voltage_sampling_frequency", fmt.Sprintf("%d", rate))
}

// SetRFBandwidth sets the RF bandwidth for both RX and TX
func (dev *AD9361) SetRFBandwidth(bw uint32) error {
	if err := dev.phy.SetAttr("in_voltage_rf_bandwidth", fmt.Sprintf("%d", bw)); err != nil {
		return err
	}
	return dev.phy.SetAttr("out_voltage_rf_bandwidth", fmt.Sprintf("%d", bw))
}

// SetGainControl sets the gain control mode
func (dev *AD9361) SetGainControl(mode string) error {
	return dev.phy.SetAttr("in_voltage0_gain_control_mode", mode)
}

// SetRXGain sets the RX gain in dB
func (dev *AD9361) SetRXGain(gain float64) error {
	return dev.phy.SetAttr("in_voltage0_hardwaregain", fmt.Sprintf("%.2f", gain))
}

// SetTXGain sets the TX gain in dB
func (dev *AD9361) SetTXGain(gain float64) error {
	return dev.phy.SetAttr("out_voltage0_hardwaregain", fmt.Sprintf("%.2f", gain))
}

// SetQuadratureTracking enables/disables quadrature tracking
func (dev *AD9361) SetQuadratureTracking(enable bool) error {
	val := "0"
	if enable {
		val = "1"
	}
	return dev.phy.SetAttr("in_voltage0_quadrature_tracking_en", val)
}

// SetDCOffsetTracking enables/disables DC offset tracking
func (dev *AD9361) SetDCOffsetTracking(enable bool) error {
	val := "0"
	if enable {
		val = "1"
	}
	return dev.phy.SetAttr("in_voltage0_dc_offset_tracking_en", val)
}

// Close releases all resources associated with the device
func (dev *AD9361) Close() error {
	dev.phy.Destroy()
	dev.tx.Destroy()
	dev.rx.Destroy()
	return nil
}

// Modulate applies the specified modulation to the input samples
func (dev *AD9361) Modulate(samples []complex64) ([]complex64, error) {

	switch dev.sdr.config.Modulation {
	case ModulationAM:
		return dev.modulateAM(samples)
	case ModulationFM:
		return dev.modulateFM(samples)
	case ModulationPSK:
		return dev.modulatePSK(samples)
	case ModulationQAM:
		return dev.modulateQAM(samples)
	default:
		return nil, fmt.Errorf("unsupported modulation type")
	}
}

func (dev *AD9361) modulateAM(samples []complex64) ([]complex64, error) {
	output := make([]complex64, len(samples))
	carrier := complex(math.Cos(dev.CarrierFrequency), math.Sin(dev.CarrierFrequency))

	for i, sample := range samples {
		// AM modulation: (1 + m*signal)*carrier
		magnitude := float32(1.0 + real(sample))
		output[i] = complex(magnitude*float32(real(carrier)), magnitude*float32(imag(carrier)))
	}

	return output, nil
}

func (dev *AD9361) modulateFM(samples []complex64) ([]complex64, error) {
	output := make([]complex64, len(samples))
	phase := 0.0

	for i, sample := range samples {
		// FM modulation: exp(j*(wc*t + m*integral(signal)))
		phase += dev.sdr.config.FMCoeff * float64(real(complex128(sample)))
		output[i] = complex64(cmplx.Rect(1.0, phase))
	}

	return output, nil
}

func (dev *AD9361) modulatePSK(samples []complex64) ([]complex64, error) {
	output := make([]complex64, len(samples))
	symbolAngle := 2 * math.Pi / float64(dev.sdr.config.PSKOrder)

	for i, sample := range samples {
		// PSK modulation: map symbols to constellation points
		symbol := int(real(sample)) % dev.sdr.config.PSKOrder
		angle := float64(symbol) * symbolAngle
		output[i] = complex64(cmplx.Rect(1.0, angle))
	}

	return output, nil
}

func (dev *AD9361) modulateQAM(samples []complex64) ([]complex64, error) {
	output := make([]complex64, len(samples))
	sqrtM := int(math.Sqrt(float64(dev.sdr.config.QAMOrder)))
	scale := 2.0 * float64(sqrtM-1)

	for i, sample := range samples {
		// QAM modulation: map symbols to square constellation
		symbol := int(real(sample)) % dev.sdr.config.QAMOrder
		I := float64(symbol%sqrtM)*2.0 - float64(sqrtM-1)
		Q := float64(symbol/sqrtM)*2.0 - float64(sqrtM-1)
		output[i] = complex64(complex(I/scale, Q/scale))
	}

	return output, nil
}

// ConfigureDDS sets up the Direct Digital Synthesis
func (dev *AD9361) ConfigureDDS() error {
	if err := dev.tx.SetAttr("out_altvoltage0_TX1_I_F1_frequency",
		fmt.Sprintf("%d", dev.DDSFrequency)); err != nil {
		return err
	}

	scale := fmt.Sprintf("%.6f", dev.DDSScale)
	if err := dev.tx.SetAttr("out_altvoltage0_TX1_I_F1_scale", scale); err != nil {
		return err
	}

	return nil
}

// SetFIRFilter configures custom FIR filter taps
func (dev *AD9361) SetFIRFilter(taps []int16) error {
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

// Print - prints device info
func (dev *AD9361) String(options *common.ConsoleOptions) string {
	var deviceList [3]*iio.Device
	deviceList[0] = dev.phy
	deviceList[1] = dev.tx
	deviceList[2] = dev.rx
	var strbuild string
	for _, device := range deviceList {
		strbuild += fmt.Sprintf("%sDEVICE:%s %s\n", common.CONSOLE_GREEN, common.CONSOLE_RESET, device.GetName())

		if options.AttrName {
			strbuild += fmt.Sprintf("	Attributes:\n")
			attrs := device.GetAttributes()
			for _, attr := range attrs {
				if attr != nil {
					strbuild += fmt.Sprintf("		%s\n", attr.GetName())

					if options.AttrVal {
						strbuild += fmt.Sprintf("		%s\n", attr.GetValue())
					}
				}
			}
		}

		strbuild += fmt.Sprintf("	Channels:\n")
		channels := device.GetChannels()
		for _, channel := range channels {
			if channel != nil {
				name := channel.DirectionalName()
				if name != "" {
					strbuild += fmt.Sprintf("		%s\n", name)
				}

			}
		}
	}
	return strbuild
}

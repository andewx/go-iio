//go:generate sudo sh ../pkg-config/macos_install.sh

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/andewx/fft"
	"github.com/andewx/go-iio/common"
	"github.com/andewx/go-iio/iio"
	"github.com/andewx/go-iio/plot"
	"github.com/andewx/go-iio/sdr"
)

func idxCheck(args []string, index int) string {
	if index+1 < len(args) {
		return args[index+1]
	}
	return ""
}

var ready bool

func main() {
	var device *sdr.SDR
	var args []string
	var err error
	var host, name string
	host = "192.168.2.1"
	name = "ad9361"
	var scan *iio.ScanBlock
	var attrval, attrname bool
	var kill chan int
	var streamRx chan int
	var fftSize int
	var bufferSize int

	bufferSize = 32768
	fftSize = 1024
	attrval = false
	attrname = false
	ready = true
	kill = make(chan int, 1)
	streamRx = make(chan int, 1)

	// Get OS Args
	args = os.Args[1:]
	if len(args) > 0 {
		// Check args for options -host, -name, -help, -attrval -attrname -chname wit switch cases
		for index, arg := range args {
			switch arg {
			case "-host", "-h":
				host = idxCheck(args, index)
			case "-name", "-n":
				name = idxCheck(args, index)
			case "-help":
				fmt.Printf("goiio [options] \n")
				fmt.Printf("Options:\n")
				fmt.Printf("  -host | -h <host> : Set the host address\n")
				fmt.Printf("  -name | -n <name> : Set the device name\n")
				fmt.Printf("  -help | -h : Show this help message\n")
				fmt.Printf("  -attr | -a : Show the attribute value\n")
				fmt.Printf("  -debug | -d : Show the debug messages\n")
				fmt.Printf("  -verbose | -v : Show the verbose messages\n")
				return
			case "-attr", "-a":
				attrval = true
				attrname = true
			case "-debug", "-d":
				common.DebugMode = true
			case "-verbose", "-v":
				common.VerboseMode = true
			}
		}
	}

	console := common.NewConsoleOptions()

	console.AttrVal = attrval
	console.AttrName = attrname
	console.Host = host
	console.Device = name

	fmt.Printf("go-iio\n")

	scan, err = iio.CreateScanBlock("ip:usb", 0)

	if err != nil {
		fmt.Printf("IIO unable to scan for available contexts\n")
		return
	}

	scan.Scan()

	fmt.Printf("...Scanning for available contexts %s\n", scan.String())
	fmt.Printf("...Testing device availability for AD9361 %s\n", iio.Version())

	device, err = sdr.NewSDR(host, common.GhZ(1.0), name, bufferSize)

	if err != nil || device == nil {
		fmt.Printf("%s Error: %s %s\n", common.CONSOLE_RED, common.CONSOLE_RESET, err.Error())
		return
	}
	fmt.Printf("%s\n", Help())
	// Configure Device and Read In Device Stream and Plot at 10FPS
	spectrum := make([]complex64, device.Params.BufferSize/2)
	go device.StreamRx(60, streamRx)
	fft.ComputeFramesOverlap(spectrum, 0.5, fftSize)
	grapher := plot.NewFFTPlotServer(device, fftSize)
	go grapher.Start(kill)
	go ListenInput(kill, device)

	// Start a 40 frames per second 1oop
	frameTime := int64((1.0 / float64(40)) * 1e9)
	for ready {
		timeSet := time.Now().UnixNano()
		timeDelta := time.Now().UnixNano() - timeSet
		for timeDelta < frameTime {
			// Sleep this thread until the next frame if this frame is not ready
			// Sleep Again by a delta amount
			time.Sleep(time.Nanosecond * time.Duration(timeDelta))
			timeDelta = time.Now().UnixNano() - timeSet
		}

		samples := device.GetStreamRxData(streamRx)
		fresult, ferr := fft.ComputeFramesOverlap(samples, 0.5, fftSize)
		if ferr != nil {
			fmt.Printf("Error line 131 ComputeFramesOverlap\n")
			kill <- 1
			ready = false
		}

		grapher.SetReal(fresult)
	}

}

// ListenInput is a function that kills the process by setting the global var ready to false
func ListenInput(signal chan int, device *sdr.SDR) {
	reader := bufio.NewReader(os.Stdin)
	quit := false

	for !quit {
		fmt.Printf(": ")
		in, err := reader.ReadString('\n')

		if err != nil {
			fmt.Printf("Error reading input\n")
			continue
		}

		args := strings.Split(in, ":")

		switch args[0] {
		case "q":
			quit = true
			break
		case "rx":
			if len(args) > 1 {
				// Convert the frequenct to a decimal value
				freqF, err := strconv.ParseFloat(strings.TrimSpace(args[1]), 64)
				if err != nil {
					fmt.Printf("Invalid frequency: %s\n", args[1])
					continue
				}

				var freq uint64
				switch strings.TrimSpace(args[2]) {
				case "ghz":
					freq = common.GhZ(freqF)
				case "mhz":
					freq = common.MhZ(freqF)
				case "khz":
					freq = common.KhZ(freqF)
				}

				device.SetRxFrequency(freq)
				device.Update()
			}
		case "tx":
			if len(args) > 1 {
				// Convert the frequenct to a decimal value
				freqF, err := strconv.ParseFloat(strings.TrimSpace(args[1]), 64)
				if err != nil {
					fmt.Printf("Invalid frequency: %s\n", args[1])
					continue
				}

				var freq uint64
				switch strings.TrimSpace(args[2]) {
				case "ghz":
					freq = common.GhZ(freqF)
				case "mhz":
					freq = common.MhZ(freqF)
				case "khz":
					freq = common.KhZ(freqF)
				}

				device.SetTxFrequency(freq)
				device.Update()
			}
		case "info":
			if len(args) > 1 {
				fmt.Print(device.Platform.String(&common.ConsoleOptions{AttrVal: true, AttrName: true}))
			}
		case "sample":

			if len(args) > 1 {
				// Convert the frequenct to a decimal value
				freqF, err := strconv.ParseFloat(strings.TrimSpace(args[1]), 64)
				if err != nil {
					fmt.Printf("Invalid frequency: %s\n", args[1])
					continue
				}

				switch strings.TrimSpace(args[2]) {
				case "ghz":
					device.SetSampleRate(common.GhZ(freqF))
				case "mhz":
					device.SetSampleRate(common.MhZ(freqF))
				case "khz":
					device.SetSampleRate(common.KhZ(freqF))
				case "hz":
					device.SetSampleRate(uint64(freqF))
				}
				device.Update()
			}
		case "bandwidth":
			if len(args) > 1 {
				// Convert the frequenct to a decimal value
				freqF, err := strconv.ParseFloat(strings.TrimSpace(args[1]), 64)
				if err != nil {
					fmt.Printf("Invalid frequency: %s\n", args[1])
					continue
				}
				switch args[2] {
				case "ghz":
					device.SetRFBandwidth(common.GhZ(freqF))
				case "mhz":
					device.SetRFBandwidth(common.MhZ(freqF))
				case "khz":
					device.SetRFBandwidth(common.KhZ(freqF))
				case "hz":
					device.SetRFBandwidth(uint64(freqF))
				}
				device.Update()
			}
		case "gain":
			if len(args) > 1 {
				gain, err := strconv.ParseFloat(args[1], 64)
				if err != nil {
					fmt.Printf("Invalid gain: %s\n", args[1])
					continue
				}
				device.Platform.SetRXGain(gain)
			}
		case "gainmode":
			if len(args) > 1 {
				gainMode := args[1]
				if err != nil {
					fmt.Printf("Invalid gain mode: %s\n", args[1])
					continue
				}
				device.Platform.SetGainControl(gainMode)
			}
		case "help":
			fmt.Printf("%s\n", Help())

		default:
			fmt.Printf("Invalid Command, type help to see available commands %s\n", args[0])
		}

	}
	signal <- 1
	ready = false
	fmt.Printf("Exiting device...\n")
}

// Help prints the command line command options
func Help() string {
	return "go-iio interpreter commands \n" +
		"Commands:\n" +
		"rx:<freq>:<ghz|mhz|khz|hz> to set the RX Frequency\n" +
		"tx:<freq>::<ghz|mhz|khz|hz> to set the TX Frequency\n" +
		"info to print the device tree\n" +
		"sample:<rate>::<ghz|mhz|khz|hz> to set the sample rate\n" +
		"gain:<value> to set the gain\n" +
		"gainmode:<value> to set the gain mode\n" +
		"read:<device>:<channel>:<attribute> to read an attribute\n" +
		"[enter]|q: to quit\n"
}

// ComplexToFloat32 converts a complex64 array to a float32 array
func ComplexToFloat32(in []complex64) [][]float32 {
	out := make([][]float32, 2)
	out[0] = make([]float32, len(in))
	out[1] = make([]float32, len(in))
	for i, v := range in {
		out[0][i] = real(v)
		out[1][i] = imag(v)
	}
	return out
}

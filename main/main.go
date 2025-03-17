//go:generate sudo sh ../pkg-config/macos_install.sh

package main

import (
	"bufio"
	"fmt"
	"os"
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
				fmt.Printf("  -attrval | -av : Show the attribute value\n")
				fmt.Printf("  -attrname | -an : Show the attribute name\n")
				fmt.Printf("  -debug | -d : Show the debug messages\n")
				fmt.Printf("  -verbose | -v : Show the verbose messages\n")
				return
			case "-attrval", "-av":
				attrval = true
			case "-attrname", "-an":
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

	fmt.Printf("Welcome...Starting SDR Test Suite\n")

	fmt.Printf("...Finding available contexts and listing\n")

	scan, err = iio.CreateScanBlock("ip:usb", 0)

	if err != nil {
		fmt.Printf("IIO unable to scan for available contexts")
		return
	}

	scan.Scan()

	fmt.Printf("...Scanning for available contexts %s\n", scan.String())
	fmt.Printf("...Testing device availability for AD9361 %s\n", iio.Version())

	device, err = sdr.NewSDR(host, name, 32678)

	if err != nil || device == nil {
		fmt.Printf("Unable to connect to Network Device %s at %s \n", host, name)
		return
	}

	device.Print(console)
	fmt.Printf("Starting FFT Plot Server\n")
	fmt.Printf("...Listening on locahost:8080\n")
	fmt.Printf("Press [enter] to exit...\n")

	// Configure Device and Read In Device Stream and Plot at 10FPS
	spectrum := make([]complex64, device.Params.BufferSize/2)
	go device.StreamRx(40, streamRx)
	fft.FFTSinglePrecision(spectrum)
	grapher := plot.NewFFTPlotServer(device, 4096)
	go grapher.Start(kill)
	go ListenInput(kill)

	// Start a 10 frames per second 1oop
	frameTime := int64((1.0 / float64(10)) * 1e9)
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
		err = fft.FFTSinglePrecision(samples)

		if err != nil {
			fmt.Printf("Error in FFT\n")
			kill <- 1
			ready = false
		}

		data := fft.PowerSpectrum(samples)

		grapher.SetReal(data)
	}

}

// ListenInput is a function that kills the process by setting the global var ready to false
func ListenInput(signal chan int) {
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	//Send message to kill
	signal <- 1
	ready = false
	fmt.Printf("Exiting device...\n")
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

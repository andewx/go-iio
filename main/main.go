//go:generate sudo sh ../pkg-config/macos_install.sh

package main

import (
	"fmt"
	"os"

	"github.com/andewx/go-iio/common"
	"github.com/andewx/go-iio/fft"
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

func main() {
	var g *sdr.SDR
	var args []string
	var err error
	var host, name string
	host = "192.168.2.1"
	name = "ad9361"
	var scan *iio.ScanBlock
	var attrval, attrname bool

	attrval = false
	attrname = false

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

	g, err = sdr.ConnectNetwork(host, name)

	if err != nil || g == nil {
		fmt.Printf("Unable to connect to Network Device %s at %s \n", host, name)
		return
	}

	g.Print(console)

	// Configure Device and Read In Device Stream and Plot
	iqData := g.GetSamples(1024, 2)
	spectrum, _ := fft.FFT64(iqData)
	grapher := plot.NewFFTPlotServer(g, spectrum)
	fmt.Printf("Starting FFT Plot Server\n")
	fmt.Printf("...Listening on locahost:8080\n")
	grapher.Start()

}

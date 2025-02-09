//go:generate sudo sh ../pkg-config/macos_install.sh

package main

import (
	"fmt"

	"github.com/andewx/go-iio/iio"
	"github.com/andewx/go-iio/sdr"
)

func main() {
	var g *sdr.SDR
	var err error
	var host, name string
	host = "192.168.2.1"
	name = "ad9361"
	var scan *iio.ScanBlock

	fmt.Printf("Welcome...Starting SDR Test Suite\n")

	fmt.Printf("...Finding available contexts and listing\n")

	scan, err = iio.CreateScanBlock("ip:usb", 0)

	if err != nil {
		fmt.Printf("IIO unable to scan for available contexts")
		return
	}

	scan.Scan()

	fmt.Printf("...Scanning for available contexts\n%s", scan.String())
	fmt.Printf("...Testing device availability for AD9361\n", iio.Version())

	g, err = sdr.ConnectNetwork(host, name)

	if err != nil || g == nil {
		fmt.Printf("Unable to connect to Network Device %s at %s \n", host, name)
		return
	}

	g.Print()

}

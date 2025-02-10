package test

import (
	"fmt"
	"testing"

	"github.com/andewx/go-iio/iio"
	"github.com/andewx/go-iio/sdr"
)

func DeviceTest(t *testing.T) {
	var g *sdr.SDR
	var err error
	var host, name string
	host = "192.168.2.0"
	name = "ad9361-phy"
	var scan *iio.ScanBlock

	fmt.Printf("Finding available contexts and listing\n")

	scan, err = iio.CreateScanBlock("ip:usb", 0)

	if err != nil {
		t.Errorf("IIO unable to scan for available contexts")
		return
	}

	scan.Scan()

	fmt.Printf("IIO LIB SCAN CONTEXTS\n%s\n", scan.String())

	fmt.Printf("go-iio %s\ntesting device availability for AD9361\n", iio.Version())

	g, err = sdr.ConnectNetwork(host, name)

	if err != nil || g == nil {
		t.Errorf("Unable to connect to Network Device %s at %s \n", host, name)
	}

	fmt.Printf("exiting...\n")

}

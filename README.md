# go-iio

#### libiio Version 0.26

**Overview**

This is a Go library for the libiio library. Note that the libiio library is a C library, so this library is a Go wrapper around the C library. Different version of the libiio library will usually require a different release branch of this repository as the C library changes significantly between versions even within the same major version.

## Installation

Recommended installation method is to use the `analogdevicesinc` github repository and build the libiio library from source. Ensure that the `libiio` repository is cloned to the same version as the `go-iio` repository.

The 0.26 version of the `libiio` repository is known to work with this repository.

Must ensure that you are working with the `v0.26` tag of the `libiio` repository.

Other `go-iio` versions may work, but this is not guaranteed. This library is not maintained. If you require an updated version of the library please fork this repository.

```
git clone https://github.com/analogdevicesinc/libiio.git
cd libiio
git checkout tags/0.26
```

Move libiio to the `/usr/local` repositories for global library linkage

`export PKG_CONFIG_PATH=/usr/local/lib/pkg-config`

- Headers are installed to `/usr/local/include/iio.h`
- Library is installed to `/usr/local/lib/libiio`


## About

This library is an API binding framework for the libiio C ABI library. Currently, it is only tested with the AD9361 radio, but it should work with any libiio compatible device.

This library is still under development and does not yet support all features of the libiio library.

If you find any issues, please file an issue on the [GitHub repository](https://github.com/andewx/go-iio/issues).

This library was developed by Brian M. Anderson for a University of Arizona project.


## Usage

### Command Line Interface 

```
goiio [options] 
Options:
  -host | -h <host> : Set the host address
  -name | -n <name> : Set the device name
  -help | -h : Show this help message
  -attrval | -av : Show the attribute value
  -attrname | -an : Show the attribute name
  -debug | -d : Show the debug messages
  -verbose | -v : Show the verbose messages
```

### Overview

Here is a usage example for connecting to a networked libiio compatible device and printing the available devices presented, device attributes, channels, and channel attributes.

```
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
```

### `iio`

The package `iio` contains the bindings for for the libiio library and works similarily to the `C` version of the library except that the C-data structures are mapped to go structs.

## Building and Running

### macOS Version 13.0 (Ventura) or later

There is an issue with the golang complier and the `libiio` library on macOS version 13.0 (Ventura) or later where the @rpath attribute is added to the binary and lost.

For now the fix is to determine the path of the `libiio` library and headers. And change the compiled binary executable in golang with the following:

First inspect the `./pkg-config/macos_install.sh` script to determine the path of the `libiio` library and headers. and then run the following commands:

```
go build ./...
cd main
go generate
./main
```

This will change the @rpath attribute to the correct path of the `libiio` library and headers.

You can verify the @rpath attribute is correct by running the following command:

```
otool -l main | grep @rpath
```

### Devices

#### AD9361

The AD9361 is a software defined radio transceiver that is compatible with the libiio library. In order to use the AD9361 in the main executable loop ensure that after the context is create you call 

Note that `sdr` in this case is the go `package` `github.com/andewx/go-iio/sdr`

```
	g, err = sdr.ConnectNetwork("192.168.2.1", "ad9361")
```

This device is typically associate with the pluto adalm and will represent a larger category of devices eventually so we will look for a more user friendly way to handle this.

To interface with another device you may be able to explore its devicetree by opening a serial connection first.

```
ls /dev/tty.*
```
```
screen /dev/tty.usbserial-XXXXX 115200
```


### Todo

- Full associate device and channels including the RX/TX I/Q Channels and implement the channel configuration with buffers
- Scan device out to JSON file for documentation and parsing
- Implement configurations with dictionary references of the device JSON data. Write attributes to device as needed.
- Perform buffer and live tranmit and recieve tests
- Implement throttling and rate limiting features
- Implement device discovery features
- Implement device configuration features
- Implement device monitoring features
- Implement device logging features
- Implement device testing features
- Implement device documentation features
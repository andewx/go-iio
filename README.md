# go-iio

#### libiio Version 0.26

**Overview**

This is a Go library for the libiio library. Note that the libiio library is a C library, so this library is a Go wrapper around the C library. Different version of the libiio library will usually require a different release branch of this repository as the C library changes significantly between versions even within the same major version.

## Installation

Recommended installation method is to use the `analogdevicesinc` github repository and build the libiio library from source. Ensure that the `libiio` repository is cloned to the same version as the `go-iio` repository.

This version of the `libiio` repository is known to work with this version of the `go-iio` repository.

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

This library tries to be a full featured wrapper around the libiio library. Currently, it is only tested with the AD9361 radio, but it should work with any libiio compatible device.

This library is still under development and does not yet support all features of the libiio library.

If you find any issues, please file an issue on the [GitHub repository](https://github.com/andewx/go-iio/issues).

This library was developed by Brian M. Anderson for a University of Arizona project.






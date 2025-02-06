# go-iio

#### libiio Version 0.26

**Overview**

This is a Go library for the libiio library. Note that the libiio library is a C library, so this library is a Go wrapper around the C library. Different version of the libiio library will usually require a different release branch of this repository as the C library changes significantly between versions even within the same major version.

## Installation

Move libiio to the `/usr/local` repositories for global library linkage

`export PKG_CONFIG_PATH=/usr/local/lib/pkg-config`

- Headers are installed to `/usr/local/include/iio.h`
- Library is installed to `/usr/local/lib/libiio`




package iio

// #cgo pkg-config: libiio
// #include <iio.h>
import (
	"C"
	"fmt"
	"unsafe"
)

/*
 Includes the Scan Functions and Context Creations Functions within the IIO Library Files
*/

func Scan(params *ContextParamsHandle, backends string) (*Scan, error) {
	scan := &Scan{}
	var err error

	c_str := unsafe.StringData(backends)

	scan.handle, err = C.iio_scan(params.handle, c_str)
	if err != nil {
		return nil, err
	}
	return &Scan{handle: scan}, nil
}

func ScanDestroy(scan *Scan) {
	C.iio_scan_destroy(scan.handle)
}

func ScanGetCount(scan *Scan) (int, error) {
	var err error
	var res C.int
	res = C.iio_scan_get_results(scan.handle)

	if res < 0 {
		err = fmt.Errorf("Error getting the number of results")
	}
	return int(res), err
}

func ScanGetDescription(scan *Scan, index int) (string, error) {
	var err error
	var res *C.char
	var len_string int
	res = C.iio_scan_get_description(scan.handle, C.int(index))
	if res == nil {
		err = fmt.Errorf("Error getting the Description")
		return "", err
	} else {
		len_string = int(C.strlen(res))
		if len_string <= 0 {
			err = fmt.Errorf("Error getting the Description")
		}
	}
	return unsafe.String(res, len_string), err
}

func ScanGetURI(scan *Scan, index int) (string, error) {
	var err error
	var res *C.char
	res = C.iio_scan_get_uri(scan.handle, C.int(index))
	var len_string int
	if res == nil {
		err = fmt.Errorf("Error getting the URI")
		return "", err
	} else {
		len_string = int(C.strlen(res))
		if len_string <= 0 {
			err = fmt.Errorf("Error getting the URI")
		}
	}
	return unsafe.String(res, len_string), err
}

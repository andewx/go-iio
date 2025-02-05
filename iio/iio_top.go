package iio

// #cgo pkg-config: libiio
// #include <iio.h>
import (
	"C"
	"fmt"
	"unsafe"
)

func StrError(err_type int) (string, error) {
	var err error
	var c_str *C.char
	var str_len int
	c_str = C.iio_strerror(C.int(err_type))

	if c_str == nil {
		err = fmt.Errorf("Error getting the error string")
		return "", err
	} else {
		str_len = int(C.strlen(c_str))
		if str_len <= 0 {
			err = fmt.Errorf("Error getting the error string")
		}
	}
	return unsafe.String(c_str, str_len), err
}

func HasBackend(params *ContextParamsHandle, backend string) (bool, error) {
	var err error
	var res C.int
	c_str := unsafe.StringData(backend)

	res = C.iio_scan_has_backend(params.handle, c_str)

	if res < 0 {
		err = fmt.Errorf("Error checking if backend exists")
	}
	return bool(res), err
}

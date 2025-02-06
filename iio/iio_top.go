package iio

// #cgo pkg-config: libiio
// #include <iio.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"unsafe"
)

/** @brief Get a string description of an error code
 * @param err The error code. Can be positive or negative.
 * @param dst A pointer to the memory area where the NULL-terminated string
 * corresponding to the error message will be stored
 * @param len The available length of the memory area, in bytes */
func StrError(err_type int) {

	if err_type == 0 {
		return
	} else {
		var c_str *C.char
		c_str = (*C.char)(C.malloc(4 * 1024))
		C.iio_strerror(C.int(err_type), c_str, 1024)
		defer C.free(unsafe.Pointer(c_str))
		fmt.Println(C.GoString(c_str))
	}
	return
}

/** @brief Check if the specified backend is available
 * @param params A pointer to a iio_context_params structure that contains
 *   context creation information; can be NULL
 * @param backend The name of the backend to query
 * @return True if the backend is available, false otherwise */
func HasBackend(backend string) bool {
	var res C._Bool
	c_str := C.CString(backend)
	res = C.iio_has_backend(c_str)
	return bool(res)
}

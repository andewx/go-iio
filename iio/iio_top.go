package iio

// #cgo pkg-config: libiio
// #include <iio.h>
// #include <stdlib.h>

import (
	"C"
	"fmt"
	"unsafe"
)

/** @brief Get a string description of an error code
 * @param err The error code. Can be positive or negative.
 * @param dst A pointer to the memory area where the NULL-terminated string
 * corresponding to the error message will be stored
 * @param len The available length of the memory area, in bytes */
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

/** @brief Check if the specified backend is available
 * @param params A pointer to a iio_context_params structure that contains
 *   context creation information; can be NULL
 * @param backend The name of the backend to query
 * @return True if the backend is available, false otherwise */
func HasBackend(params *ContextParamsHandle, backend string) (bool, error) {
	var err error
	var res C._Bool
	c_str := unsafe.StringData(backend)

	res = C.iio_has_backend(params.handle, c_str)

	if res == false {
		err = fmt.Errorf("Error checking if backend exists")
	}
	return bool(res), err
}

/** @brief Get the number of available built-in backends
 * @return The number of available built-in backends */
func GetBuiltinBackendsCount(index int) (int, error) {
	var err error
	var res C.int
	res = C.iio_get_builtin_backends_count(C.int(index))
	if res < 0 {
		err = fmt.Errorf("Error getting the number of results")
	}
	return int(res), err
}

/** @brief Retrieve the name of a given built-in backend
 * @param index The index corresponding to the backend
 * @return On success, a pointer to a static NULL-terminated string
 * @return If the index is invalid, NULL is returned */
func GetBuiltinBackend(index int) (string, error) {
	var err error
	var res *C.char
	var len_string int
	res = C.iio_get_builtin_backend(C.int(index))
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

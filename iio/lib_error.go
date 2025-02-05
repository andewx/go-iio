package iio

// #cgo pkg-config: libiio
// #include <iio.h>
// #include <errno.h>
// int get_errno() { return errno; }
import "C"
import (
	"fmt"
	"syscall"
)

func getLastError() error {
	var errno C.int
	errno = C.get_errno()
	return syscall.Errno(errno)
}

func getError(code C.ssize_t) error {
	if code < 0 {
		return syscall.Errno(-code)
	}
	return fmt.Errorf("unknown error: %d", code)
}

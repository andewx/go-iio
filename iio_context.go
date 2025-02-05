package iio

// #cgo pkg-config: libiio
// #include <iio.h>
import (
	"C"
)
import (
	"fmt"
	"unsafe"
)

/** @} */ /* ------------------------------------------------------------------*/
/* ------------------------- Context functions -------------------------------*/
/** @defgroup Context Context
 * @{
 * @struct iio_context
 * @brief Contains the representation of an IIO context */

/** @brief Create a context from a URI description
 * @param params A pointer to a iio_context_params structure that contains
 *   context creation information; can be NULL
 * @param uri a URI describing the context location. If NULL, the backend will
 *   be created using the URI string present in the IIOD_REMOTE environment
 *   variable, or if not set, a local backend is created.
 * @return On success, a pointer to a iio_context structure
 * @return On failure, a pointer-encoded error is returned
 *
 * <b>NOTE:</b> The following URIs are supported based on compile time backend
 * support:
 * - Local backend, "local:"\n
 *   Does not have an address part. For example <i>"local:"</i>
 * - XML backend, "xml:"\n Requires a path to the XML file for the address part.
 *   For example <i>"xml:/home/user/file.xml"</i>
 * - Network backend, "ip:"\n Requires a hostname, IPv4, or IPv6 to connect to
 *   a specific running IIO Daemon or no address part for automatic discovery
 *   when library is compiled with ZeroConf support. For example
 *   <i>"ip:192.168.2.1"</i>, <b>or</b> <i>"ip:localhost"</i>, <b>or</b> <i>"ip:"</i>
 *   <b>or</b> <i>"ip:plutosdr.local"</i>. To support alternative port numbers the
 *   standard <i>ip:host:port</i> format is used. A special format is required as
 *   defined in RFC2732 for IPv6 literal hostnames, (adding '[]' around the host)
 *   to use a <i>ip:[x:x:x:x:x:x:x:x]:port</i> format.
 *   Valid examples would be:
 *     - ip:                                               Any host on default port
 *     - ip::40000                                         Any host on port 40000
 *     - ip:analog.local                                   Default port
 *     - ip:brain.local:40000                              Port 40000
 *     - ip:192.168.1.119                                  Default Port
 *     - ip:192.168.1.119:40000                            Port 40000
 *     - ip:2601:190:400:da:47b3:55ab:3914:bff1            Default Port
 *     - ip:[2601:190:400:da:9a90:96ff:feb5:acaa]:40000    Port 40000
 *     - ip:fe80::f14d:3728:501e:1f94%eth0                 Link-local through eth0, default port
 *     - ip:[fe80::f14d:3728:501e:1f94%eth0]:40000         Link-local through eth0, port 40000
 * - USB backend, "usb:"\n When more than one usb device is attached, requires
 *   bus, address, and interface parts separated with a dot. For example
 *   <i>"usb:3.32.5"</i>. Where there is only one USB device attached, the shorthand
 *   <i>"usb:"</i> can be used.
 * - Serial backend, "serial:"\n Requires:
 *     - a port (/dev/ttyUSB0),
 *     - baud_rate (default <b>115200</b>)
 *     - serial port configuration
 *        - data bits (5 6 7 <b>8</b> 9)
 *        - parity ('<b>n</b>' none, 'o' odd, 'e' even, 'm' mark, 's' space)
 *        - stop bits (<b>1</b> 2)
 *        - flow control ('<b>\0</b>' none, 'x' Xon Xoff, 'r' RTSCTS, 'd' DTRDSR)
 *
 *  For example <i>"serial:/dev/ttyUSB0,115200"</i> <b>or</b> <i>"serial:/dev/ttyUSB0,115200,8n1"</i>*/

func CreateContext(params *ContextParamsHandle, uri string) (*Context, error) {
	var err error
	var res *C.struct_iio_context
	c_str := unsafe.StringData(uri)
	res = C.iio_create_context(params.handle, c_str) // TODO RETURNS CODED POINTER ERROR
	if res == nil {
		err = fmt.Errorf("Error creating the context")
	}
	return &Context{handle: res}, err
}

/** @brief Destroy the given context
 * @param ctx A pointer to an iio_context structure
 *
 * <b>NOTE:</b> After that function, the iio_context pointer shall be invalid. */
func DestroyContext(ctx *Context) {
	C.iio_context_destroy(ctx.handle)
}

/** @brief Get the major number of the library version
 * @param ctx Optional pointer to an iio_context structure
 * @return The major number
 *
 * NOTE: If ctx is non-null, it will return the major version of the remote
 * library, if running remotely. */
func GetVersionMajor(ctx *Context) int {
	return int(C.iio_context_get_version_major(ctx.handle))
}

/** @brief Get the minor number of the library version
 * @param ctx Optional pointer to an iio_context structure
 * @return The minor number
 *
 * NOTE: If ctx is non-null, it will return the minor version of the remote
 * library, if running remotely. */
func GetVersionMinor(ctx *Context) int {
	return int(C.iio_context_get_version_minor(ctx.handle))
}

/** @brief Get the git hash string of the library version
 * @param ctx Optional pointer to an iio_context structure
 * @return A NULL-terminated string that contains the git tag or hash
 *
 * NOTE: If ctx is non-null, it will return the git tag or hash of the remote
 * library, if running remotely. */
func GetVersionTag(ctx *Context) (string, error) {
	var err error
	c_str := C.iio_context_get_version_tag(ctx.handle)
	if c_str == nil {
		err = fmt.Errorf("Error getting the error string")
	} else {
		str_len := int(C.strlen(c_str))
		if str_len <= 0 {
			err = fmt.Errorf("Error getting the error string")
		}
	}
	return unsafe.String(c_str), err
}

/** @brief Obtain a XML representation of the given context
 * @param ctx A pointer to an iio_context structure
 * @return On success, an allocated string. Must be deallocated with free().
 * @return On failure, a pointer-encoded error is returned */

func GetContextXML(ctx *Context) (string, error) {
	//WARNING UNSAFE Function with free not implemented yet
	var err error
	c_str := C.iio_context_get_xml(ctx.handle) //Allocated string must be deallocate wtih free not implemented yet
	if c_str == nil {
		err = fmt.Errorf("Error getting the error string")
	}
	return unsafe.String(c_str), err
}

/** @brief Get the name of the given context
 * @param ctx A pointer to an iio_context structure
 * @return A pointer to a static NULL-terminated string
 *
 * <b>NOTE:</b>The returned string will be <b><i>local</i></b>,
 * <b><i>xml</i></b> or <b><i>network</i></b> when the context has been
 * created with the local, xml and network backends respectively.*/
func GetContextName(ctx *Context) string {
	var err error
	c_str := C.iio_context_get_name(ctx.handle)
	if c_str == nil {
		err = fmt.Errorf("Error getting the error string")
	} else {
		str_len := int(C.strlen(c_str))
		if str_len <= 0 {
			err = fmt.Errorf("Error getting the error string")
		}
	}
	return unsafe.String(c_str), err
}

/** @brief Get a description of the given context
 * @param ctx A pointer to an iio_context structure
 * @return A pointer to a static NULL-terminated string
 *
 * <b>NOTE:</b>The returned string will contain human-readable information about
 * the current context. */
func GetContextDescription(ctx *Context) string {
	var err error
	c_str := C.iio_context_get_description(ctx.handle)
	if c_str == nil {
		err = fmt.Errorf("Error getting the error string")
	} else {
		str_len := int(C.strlen(c_str))
		if str_len <= 0 {
			err = fmt.Errorf("Error getting the error string")
		}
	}
	return unsafe.String(c_str), err
}

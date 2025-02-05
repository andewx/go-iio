package iio

// #cgo pkg-config: libiio
// #include <iio.h>
// #include <stdlib.h>
import "C"
import (
	"unsafe"
)

// Context represents an IIO context
type Context struct {
	handle *C.struct_iio_context
}

// CreateContext creates the default IIO context
func CreateContext() (*Context, error) {
	handle := C.iio_create_default_context()
	if handle == nil {
		return nil, getLastError()
	}
	return &Context{handle: handle}, nil
}

// CreateNetworkContext creates a network context from the specified hostname
func CreateNetworkContext(hostname string) (*Context, error) {
	cHostname := C.CString(hostname)
	defer C.free(unsafe.Pointer(cHostname))
	handle := C.iio_create_network_context(cHostname)
	if handle == nil {
		return nil, getLastError()
	}
	return &Context{handle: handle}, nil
}

// CreateXMLContext creates a context from an XML file
func CreateXMLContext(xmlFile string) (*Context, error) {
	cXmlFile := C.CString(xmlFile)
	defer C.free(unsafe.Pointer(cXmlFile))
	handle := C.iio_create_xml_context(cXmlFile)
	if handle == nil {
		return nil, getLastError()
	}
	return &Context{handle: handle}, nil
}

// CreateURIContext creates a context from an URI description
func CreateURIContext(uri string) (*Context, error) {
	cUri := C.CString(uri)
	defer C.free(unsafe.Pointer(cUri))
	handle := C.iio_create_context_from_uri(cUri)
	if handle == nil {
		return nil, getLastError()
	}
	return &Context{handle: handle}, nil
}

// GetVersion returns the backend version information
func (ctx *Context) GetVersion() (major, minor, git_tag string) {
	var cMajor, cMinor, cGitTag *C.char
	C.iio_context_get_version(ctx.handle, &cMajor, &cMinor, &cGitTag)
	return C.GoString(cMajor), C.GoString(cMinor), C.GoString(cGitTag)
}

// GetName returns the name of the backend
func (ctx *Context) GetName() string {
	return C.GoString(C.iio_context_get_name(ctx.handle))
}

// GetDescription returns the description of the context
func (ctx *Context) GetDescription() string {
	return C.GoString(C.iio_context_get_description(ctx.handle))
}

// GetXMLAttr returns the value of an XML attribute
func (ctx *Context) GetXMLAttr(name string) (string, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var value *C.char
	ret := C.iio_context_get_attr_value(ctx.handle, cName, &value)
	if ret < 0 {
		return "", getError(ret)
	}
	return C.GoString(value), nil
}

// GetDevicesCount returns the number of devices in the context
func (ctx *Context) GetDevicesCount() uint {
	return uint(C.iio_context_get_devices_count(ctx.handle))
}

// GetDevice returns the device at the specified index
func (ctx *Context) GetDevice(index uint) (*Device, error) {
	handle := C.iio_context_get_device(ctx.handle, C.uint(index))
	if handle == nil {
		return nil, getLastError()
	}
	return &Device{handle: handle}, nil
}

// FindDevice finds a device by its name or ID
func (ctx *Context) FindDevice(name string) (*Device, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	handle := C.iio_context_find_device(ctx.handle, cName)
	if handle == nil {
		return nil, getLastError()
	}
	return &Device{handle: handle}, nil
}

// Close destroys the context
func (ctx *Context) Close() error {
	if ctx.handle != nil {
		C.iio_context_destroy(ctx.handle)
		ctx.handle = nil
	}
	return nil
}

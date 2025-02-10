package iio

// #cgo pkg-config: libiio
// #include <iio.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/andewx/go-iio/common"
)

// Context represents an IIO context
type Context struct {
	handle *C.struct_iio_context
	attrs  []*ContextAttribute
}

// CreateContext creates the default IIO context
func CreateContext() (*Context, error) {
	var err error
	handle := C.iio_create_default_context()
	if handle == nil {
		return nil, getLastError()
	}
	ctx := &Context{handle: handle, attrs: nil}

	ctx.attrs = GetContextAttributes(handle)

	if ctx.attrs == nil {
		err = fmt.Errorf("Default Context attributes not available\n")
	}

	return ctx, err
}

// CreateNetworkContext creates a network context from the specified hostname
func CreateNetworkContext(hostname string) (*Context, error) {
	common.PrintDebug("File:iio_context.go // Line:38 // Function:CreateNetworkContext // Creating network context", hostname)
	var err error
	cHostname := C.CString(hostname)
	defer C.free(unsafe.Pointer(cHostname))
	common.PrintDebug("File:iio_context.go // Line:44 //Creating Network Device with libiio function iio_create_network_context() ->", cHostname)
	handle := C.iio_create_network_context(cHostname)
	if handle == nil {
		return nil, getLastError()
	}
	common.PrintValid("Network Device created successfully")
	ctx := &Context{handle: handle, attrs: nil}
	common.PrintDebug("File:iio_context.go // Line:55 // Getting Context Attributes", ctx)
	ctx.attrs = GetContextAttributes(handle)
	common.PrintValid("Network Context attributes retrieved successfully")
	if ctx.attrs == nil {
		err = fmt.Errorf("Network Context attributes not available\n")
	}

	common.PrintDebug("File:iio_context.go // Line:59 // Returning Context", ctx)

	return ctx, err
}

// CreateXMLContext creates a context from an XML file
func CreateXMLContext(xmlFile string) (*Context, error) {
	var err error
	cXmlFile := C.CString(xmlFile)
	defer C.free(unsafe.Pointer(cXmlFile))
	handle := C.iio_create_xml_context(cXmlFile)
	if handle == nil {
		return nil, getLastError()
	}

	ctx := &Context{handle: handle, attrs: nil}

	ctx.attrs = GetContextAttributes(handle)

	if ctx.attrs == nil {
		err = fmt.Errorf("Default Context attributes not available\n")
	}

	return ctx, err
}

// CreateURIContext creates a context from an URI description
func CreateURIContext(uri string) (*Context, error) {
	var err error
	cUri := C.CString(uri)
	defer C.free(unsafe.Pointer(cUri))
	handle := C.iio_create_context_from_uri(cUri)
	if handle == nil {
		return nil, getLastError()
	}

	ctx := &Context{handle: handle, attrs: nil}

	ctx.attrs = GetContextAttributes(handle)

	if ctx.attrs == nil {
		err = fmt.Errorf("Default Context attributes not available\n")
	}

	return ctx, err
}

// GetVersion returns the backend version information
func (ctx *Context) GetVersion() (major int, minor int, git_tag string) {
	var cMajor, cMinor *C.uint
	var cGitTag [9]C.char
	cGitTag[8] = '\000'
	C.iio_context_get_version(ctx.handle, cMajor, cMinor, &cGitTag[0])
	return int(*cMajor), int(*cMinor), C.GoString(&cGitTag[0])
}

// GetName returns the name of the backend
func (ctx *Context) GetName() string {
	return C.GoString(C.iio_context_get_name(ctx.handle))
}

// GetDescription returns the description of the context
func (ctx *Context) GetDescription() string {
	return C.GoString(C.iio_context_get_description(ctx.handle))
}

func (ctx *Context) GetAttrsCount() int {
	return int(C.iio_context_get_attrs_count(ctx.handle))
}

func (ctx *Context) GetAttr(index int) (string, string, error) {
	var name, value *C.char
	ret := C.iio_context_get_attr(ctx.handle, C.uint(index), &name, &value)
	if ret < 0 {
		return "", "", fmt.Errorf("iio: attribute %d not found", index)
	}
	return C.GoString(name), C.GoString(value), nil
}

// GetXMLAttr returns the value of an XML attribute
func (ctx *Context) GetXMLAttr(name string) (string, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var value *C.char
	ret := C.iio_context_get_attr_value(ctx.handle, cName)
	if ret == nil {
		fmt.Printf("iio: attribute %s not found", name)
	}
	return C.GoString(value), nil
}

// GetDevicesCount returns the number of devices in the context
func (ctx *Context) GetDevicesCount() uint {
	return uint(C.iio_context_get_devices_count(ctx.handle))
}

func (ctx *Context) GetDevices() ([]*Device, error) {
	count := ctx.GetDevicesCount()
	devices := make([]*Device, count)
	for i := uint(0); i < count; i++ {
		device, err := ctx.GetDevice(i)
		if err != nil {
			return nil, err
		}
		devices[i] = device
	}
	return devices, nil
}

// GetDevice returns the device at the specified index
func (ctx *Context) GetDevice(index uint) (*Device, error) {
	handle := C.iio_context_get_device(ctx.handle, C.uint(index))
	if handle == nil {
		return nil, getLastError()
	}
	return NewDevice(handle), nil
}

// FindDevice finds a device by its name or ID
func (ctx *Context) FindDevice(name string) (*Device, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	handle := C.iio_context_find_device(ctx.handle, cName)
	if handle == nil {
		return nil, getLastError()
	}
	dev := NewDevice(handle)
	if dev == nil {
		return nil, fmt.Errorf("%s device couldn't be initiated", name)
	}
	common.PrintDebug(fmt.Sprintf("file iio_context.go::FindDevice|%s", dev._name), dev)
	return dev, nil
}

// Close destroys the context
func (ctx *Context) Close() error {
	if ctx.handle != nil {
		C.iio_context_destroy(ctx.handle)
		ctx.handle = nil
	}
	return nil
}

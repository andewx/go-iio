package iio

/*
#include <iio.h>
#include <errno.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

// ScanContext represents a context for scanning available IIO devices
type ScanContext struct {
	ctx *C.struct_iio_scan_context
}

// ContextInfo contains information about a discovered IIO context
type ContextInfo struct {
	info *C.struct_iio_context_info
}

// ScanBlock represents a block for scanning IIO contexts
type ScanBlock struct {
	blk *C.struct_iio_scan_block
}

// CreateScanContext creates a new scan context for discovering IIO devices
// backend is a comma-separated list of backends to use for scanning (e.g. "local,usb,ip")
// flags should be set to 0 for now as it's unused in the C API
func CreateScanContext(backend string, flags uint) (*ScanContext, error) {
	cBackend := C.CString(backend)
	defer C.free(unsafe.Pointer(cBackend))

	ctx := C.iio_create_scan_context(cBackend, C.uint(flags))
	if ctx == nil {
		return nil, errors.New("failed to create scan context")
	}

	return &ScanContext{ctx: ctx}, nil
}

// Destroy frees the resources associated with the scan context
func (s *ScanContext) Destroy() {
	if s.ctx != nil {
		C.iio_scan_context_destroy(s.ctx)
		s.ctx = nil
	}
}

// GetInfoList enumerates available contexts and returns a list of context info
func (s *ScanContext) GetInfoList() ([]*ContextInfo, error) {
	var infoPtr **C.struct_iio_context_info
	count := C.iio_scan_context_get_info_list(s.ctx, &infoPtr)

	if count < 0 {
		return nil, errors.New("failed to get context info list")
	}

	// Create slice to hold the results
	result := make([]*ContextInfo, count)

	// Convert C array to Go slice
	slice := unsafe.Slice(infoPtr, count)
	for i := 0; i < int(count); i++ {
		result[i] = &ContextInfo{info: slice[i]}
	}

	return result, nil
}

// FreeInfoList frees the memory allocated for the context info list
func FreeInfoList(infoList []*ContextInfo) {
	if len(infoList) == 0 {
		return
	}

	// Get pointer to first element
	infoPtr := (**C.struct_iio_context_info)(unsafe.Pointer(&infoList[0].info))
	C.iio_context_info_list_free(infoPtr)
}

// GetDescription returns the description of a discovered context
func (ci *ContextInfo) GetDescription() string {
	return C.GoString(C.iio_context_info_get_description(ci.info))
}

// GetURI returns the URI of a discovered context
func (ci *ContextInfo) GetURI() string {
	return C.GoString(C.iio_context_info_get_uri(ci.info))
}

// CreateScanBlock creates a new scan block for discovering IIO contexts
func CreateScanBlock(backend string, flags uint) (*ScanBlock, error) {
	var cBackend *C.char
	if backend != "" {
		cBackend = C.CString(backend)
		defer C.free(unsafe.Pointer(cBackend))
	}

	blk := C.iio_create_scan_block(cBackend, C.uint(flags))
	if blk == nil {
		return nil, errors.New("failed to create scan block")
	}

	return &ScanBlock{blk: blk}, nil
}

// Destroy frees the resources associated with the scan block
func (sb *ScanBlock) Destroy() {
	if sb.blk != nil {
		C.iio_scan_block_destroy(sb.blk)
		sb.blk = nil
	}
}

// Scan enumerates available contexts using the scan block
func (sb *ScanBlock) Scan() (int, error) {
	count := C.iio_scan_block_scan(sb.blk)
	if count < 0 {
		return 0, errors.New("failed to scan for contexts")
	}
	return int(count), nil
}

// GetInfo retrieves context info for a particular index
func (sb *ScanBlock) GetInfo(index uint) (*ContextInfo, error) {
	info := C.iio_scan_block_get_info(sb.blk, C.uint(index))
	if info == nil {
		return nil, errors.New("failed to get context info")
	}
	return &ContextInfo{info: info}, nil
}

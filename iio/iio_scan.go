package iio

// #cgo pkg-config: libiio

// #include <iio.h>
// #include <errno.h>
// #include <stdlib.h>
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

// ContextInfo contains information about a discovered IIO context
type ContextInfo struct {
	handle *C.struct_iio_context_info
}

// ScanBlock Encapsulates the Scan Context operations and allows us to associate a block of context infos
type ScanBlock struct {
	blk   *C.struct_iio_scan_block
	info  []*ContextInfo
	count int
}

// GetDescription returns the description of a discovered context
func (ci *ContextInfo) GetDescription() string {
	return C.GoString(C.iio_context_info_get_description(ci.handle))
}

// GetURI returns the URI of a discovered context
func (ci *ContextInfo) GetURI() string {
	return C.GoString(C.iio_context_info_get_uri(ci.handle))
}

/* ----------------
Scan Block
----------------- */

// CreateScanBlock creates a new scan block for discovering IIO contexts
func CreateScanBlock(backend string, flags uint) (*ScanBlock, error) {
	var cBackend *C.char
	var err error
	var count int
	if backend != "" {
		cBackend = C.CString(backend)
		defer C.free(unsafe.Pointer(cBackend))
	}

	blk := C.iio_create_scan_block(cBackend, C.uint(flags))
	if blk == nil {
		return nil, errors.New("failed to create scan block")
	}

	scn := &ScanBlock{blk: blk}
	count, err = scn.Scan()
	scn.count = count

	return &ScanBlock{blk: blk}, err
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
	return &ContextInfo{handle: info}, nil
}

// StringInfo returns the given context info description for a context
func (sb *ScanBlock) StringInfo(index int) (string, error) {
	var err error
	var info *ContextInfo
	var description string
	var uri string

	info, err = sb.GetInfo(uint(index))

	if err != nil {
		return fmt.Sprintf("Context Error at scan index %d\n", index), err
	}
	uri = info.GetURI()
	description = info.GetDescription()

	if uri == "" {
		uri = "error"
	}

	if description == "" {
		description = "error"
	}

	return fmt.Sprintf("Context Found: [%s]\n@ %s\n", uri, description), err
}

func (sb *ScanBlock) String() string {
	var strbuild string
	for i := 0; i < sb.count; i++ {
		nstring, err := sb.StringInfo(i)
		if err != nil {
			return strbuild
		}
		strbuild = fmt.Sprintf("%s\n", strbuild+nstring)
	}

	return strbuild
}

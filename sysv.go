// +build linux darwin

package shm

// #include "shm.h"
// #include <stdlib.h>
// #cgo CFLAGS: -Wall
import "C"
import (
	"fmt"
	"io"
	"os"
	"unsafe"
)

// flags for shmget
const (
	IPCCreate    = C.IPC_CREAT
	IPCExclusive = C.IPC_EXCL
	// ShmHugeTable = C.SHM_HUGETLB
	// ShmHuge2MB   = C.SHM_HUGE_2MB
	// ShmHuge1GB   = C.SHM_HUGE_1GB
	// ShmNoReserve = C.SHM_NORESERVE
)

type systemVSharedMemory struct {
	id     int
	addr   unsafe.Pointer
	size   int64
	offset int64
}

func newSystemVSharedMemory(size int64, flag int, perm os.FileMode, path string, projectID int) (*systemVSharedMemory, error) {
	pathname := C.CString(path)
	defer C.free(unsafe.Pointer(pathname))

	var addr unsafe.Pointer

	id := C.sysv_shm_create(pathname, C.int(projectID), C.size_t(size), C.int(flag), C.int(perm), &addr)
	if id == -1 {
		return nil, fmt.Errorf("create System V shared memory failed")
	}

	return &systemVSharedMemory{
		id:     int(id),
		addr:   addr,
		size:   size,
		offset: 0,
	}, nil
}

func (s *systemVSharedMemory) Read(p []byte) (int, error) {
	if s.offset >= s.size {
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}

	count := int64(len(p))
	if s.size-s.offset < count {
		count = s.size - s.offset
	}
	buffer := C.malloc(C.size_t(count))
	if buffer == nil {
		return 0, fmt.Errorf("malloc failed")
	}
	defer C.free(buffer)

	bytesRead := C.sysv_shm_read(buffer, s.addr, C.int(s.offset), C.size_t(count))
	if bytesRead == -1 {
		return 0, fmt.Errorf("read failed")
	}

	copy(p, C.GoBytes(buffer, C.int(bytesRead)))
	s.offset += int64(bytesRead)
	return int(bytesRead), nil
}

func (s *systemVSharedMemory) Write(p []byte) (n int, err error) {
	if s.offset >= s.size {
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.ErrShortWrite
	}

	count := int64(len(p))
	if s.size-s.offset < count {
		count = s.size - s.offset
	}
	bytesWrite := C.sysv_shm_write(s.addr, C.int(s.offset), unsafe.Pointer(&p[0]), C.size_t(count))
	if bytesWrite == -1 {
		return 0, fmt.Errorf("write failed")
	}

	s.offset += int64(bytesWrite)
	if int(bytesWrite) < len(p) {
		err = io.ErrShortWrite
	}
	return int(bytesWrite), err
}

func (s *systemVSharedMemory) Seek(offset int64, whence int) (int64, error) {
	var newOffset int64
	switch whence {
	case io.SeekStart:
		newOffset = offset
	case io.SeekCurrent:
		newOffset = s.offset + offset
	case io.SeekEnd:
		newOffset = s.size + offset
	default:
		return 0, fmt.Errorf("invalid whence value")
	}

	if newOffset < 0 || newOffset >= s.size {
		return 0, fmt.Errorf("offset out of range")
	}

	s.offset = newOffset
	return s.offset, nil
}

func (s *systemVSharedMemory) Close() error {
	code := C.sysv_shm_remove(C.int(s.id), s.addr)
	if code == -1 {
		return fmt.Errorf("remove System V shared memory failed")
	}
	return nil
}

func (s *systemVSharedMemory) Fd() uintptr {
	return ^(uintptr(0))
}

func (s *systemVSharedMemory) Raw() unsafe.Pointer {
	return s.addr
}

func (s *systemVSharedMemory) Size() int64 {
	return s.size
}

func (s *systemVSharedMemory) Offset() int64 {
	return s.offset
}

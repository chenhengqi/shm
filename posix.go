// +build linux darwin

package shm

// #include "shm.h"
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"io"
	"os"
	"unsafe"
)

const (
	namePrefix = "/shm"
)

type posixSharedMemory struct {
	name   string
	fd     int
	addr   unsafe.Pointer
	size   int64
	offset int64
}

func newPosixSharedMemory(size int64, flag int, perm os.FileMode) (*posixSharedMemory, error) {
	name := fmt.Sprintf("%s_%s", namePrefix, randomName())
	shmName := C.CString(name)
	defer C.free(unsafe.Pointer(shmName))

	var addr unsafe.Pointer

	fd := C.posix_create_shm(shmName, C.int(flag), C.mode_t(perm), C.off_t(size), &addr)
	if fd == -1 {
		return nil, fmt.Errorf("create POSIX shared memory failed")
	}

	return &posixSharedMemory{
		name:   name,
		fd:     int(fd),
		addr:   addr,
		size:   size,
		offset: 0,
	}, nil
}

func (s *posixSharedMemory) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (s *posixSharedMemory) Write(p []byte) (n int, err error) {
	return 0, nil
}

// see lseek(2)
func (s *posixSharedMemory) Seek(offset int64, whence int) (int64, error) {
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

func (s *posixSharedMemory) Close() error {
	shmName := C.CString(s.name)
	defer C.free(unsafe.Pointer(shmName))

	code := C.posix_destroy_shm(shmName, s.addr, C.size_t(size))
	if code == -1 {
		return fmt.Errorf("destroy POSIX shared memory failed")
	}
	return nil
}

func (s *posixSharedMemory) Fd() uintptr {
	return uintptr(s.fd)
}

func (s *posixSharedMemory) Raw() unsafe.Pointer {
	return s.addr
}

func (s *posixSharedMemory) Size() int64 {
	return s.size
}

func (s *posixSharedMemory) Offset() int64 {
	return s.offset
}

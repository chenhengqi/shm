// +build linux darwin

package shm

// #include "shm.h"
// #include <stdlib.h>
// #cgo LDFLAGS: -lrt
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

const (
	namePrefix = "/shm"
)

type posixSharedMemory struct {
	fd     int
	offset int64
	addr   unsafe.Pointer
}

func newPosixSharedMemory(size int64, flag int, perm os.FileMode) (*posixSharedMemory, error) {
	shmName := fmt.Sprintf("%s_%s", namePrefix, randomName())
	name := C.CString(shmName)
	defer C.free(unsafe.Pointer(name))

	var addr unsafe.Pointer

	fd := C.posix_create_shm(name, C.int(flag), C.mode_t(perm), C.off_t(size), &addr)
	if fd == -1 {
		return nil, fmt.Errorf("create POSIX shared memory failed")
	}

	return &posixSharedMemory{
		fd:     int(fd),
		offset: 0,
		addr:   addr,
	}, nil
}

func (s *posixSharedMemory) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (s *posixSharedMemory) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (s *posixSharedMemory) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (s *posixSharedMemory) Close() error {
	return nil
}

func (s *posixSharedMemory) Raw() unsafe.Pointer {
	return s.addr
}

func (s *posixSharedMemory) Offset() int64 {
	return s.offset
}

func (s *posixSharedMemory) Rewind() {
	s.offset = 0
}

func (s *posixSharedMemory) Fd() uintptr {
	return uintptr(s.fd)
}

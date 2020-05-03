// +build linux darwin

package shm

// #include "shm.h"
// #include <stdlib.h>
// #cgo CFLAGS: -Wall
// #cgo linux LDFLAGS: -lrt
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

	fd := C.posix_shm_create(shmName, C.int(flag), C.mode_t(perm), C.off_t(size), &addr)
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

func (s *posixSharedMemory) Read(p []byte) (int, error) {
	if s.offset >= s.size {
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}

	bytesToRead := int64(len(p))
	if s.size-s.offset < bytesToRead {
		bytesToRead = s.size - s.offset
	}

	buffer := C.malloc(C.size_t(bytesToRead))
	if buffer == nil {
		return 0, fmt.Errorf("malloc failed")
	}
	defer C.free(buffer)

	bytesRead := C.posix_shm_read(C.int(s.fd), buffer, C.size_t(bytesToRead))
	if bytesRead == -1 {
		return 0, fmt.Errorf("read failed")
	}

	copy(p, C.GoBytes(buffer, C.int(bytesRead)))
	s.offset += int64(bytesRead)
	return int(bytesRead), nil
}

func (s *posixSharedMemory) Write(p []byte) (n int, err error) {
	if s.offset >= s.size {
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.ErrShortWrite
	}

	bytesToWrite := int64(len(p))
	if s.size-s.offset < bytesToWrite {
		bytesToWrite = s.size - s.offset
	}

	bytesWrite := C.posix_shm_write(C.int(s.fd), unsafe.Pointer(&p[0]), C.size_t(bytesToWrite))
	if bytesWrite == -1 {
		return 0, fmt.Errorf("write failed")
	}

	n = int(bytesWrite)
	s.offset += int64(bytesWrite)
	if int(bytesWrite) < len(p) {
		err = io.ErrShortWrite
	}
	return n, err
}

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

	seekOffset := C.posix_shm_seek(C.int(s.fd), C.off_t(newOffset), C.int(whence))
	if seekOffset == -1 {
		return 0, fmt.Errorf("lseek failed")
	}

	s.offset = newOffset
	return s.offset, nil
}

func (s *posixSharedMemory) Close() error {
	shmName := C.CString(s.name)
	defer C.free(unsafe.Pointer(shmName))

	code := C.posix_shm_remove(shmName, s.addr, C.size_t(s.size))
	if code == -1 {
		return fmt.Errorf("remove POSIX shared memory failed")
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

package shm

import (
	"fmt"
	"io"
	"os"
	"unsafe"
)

// SharedMemory is the interface that wraps the common shared memory operations.
type SharedMemory interface {
	io.ReadWriteCloser
	io.Seeker
	Fd() uintptr
	Raw() unsafe.Pointer
	Size() int64
	Offset() int64
}

// NewPosix creates POSIX shared memory object
func NewPosix(size int64, flag int, perm os.FileMode) (SharedMemory, error) {
	return newPosixSharedMemory(size, flag, perm)
}

// NewSystemV allocates System V shared memory segment
func NewSystemV(size int64, flag int, perm os.FileMode, pathname string, projectID int) (SharedMemory, error) {
	if pathname == "" {
		return nil, fmt.Errorf("pathname must refer to an existing, accessible file")
	}
	if projectID == 0 {
		return nil, fmt.Errorf("projectID must be non-zero")
	}
	return newSystemVSharedMemory(size, flag, perm, pathname, projectID)
}

// NewSystemVPrivate allocates System V shared memory segment with IPC_PRIVATE
func NewSystemVPrivate(size int64, flag int, perm os.FileMode) (SharedMemory, error) {
	return newSystemVSharedMemory(size, flag, perm, "", 0)
}

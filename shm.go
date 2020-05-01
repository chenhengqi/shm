package shm

import (
	"io"
	"os"
	"unsafe"
)

// SharedMemory is the interface that wraps the common shared memory operations.
type SharedMemory interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
	Raw() unsafe.Pointer
	Offset() int64
	Rewind()
	Fd() uintptr
}

// NewPosix creates POSIX shared memory object
func NewPosix(size int64, flag int, perm os.FileMode) (SharedMemory, error) {
	return newPosixSharedMemory(size, flag, perm)
}

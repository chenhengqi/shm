package shm

import (
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

package shm

import (
	"os"
	"testing"
)

const (
	size = 1024 * 1024
	flag = os.O_RDWR | os.O_CREATE | os.O_EXCL
	perm = 0666
)

func TestCreatePosixSharedMemory(t *testing.T) {
	_, err := NewPosix(size, flag, perm)
	if err != nil {
		t.Fatal(err)
	}
}

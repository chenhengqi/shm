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
	mem, err := NewPosix(size, flag, perm)
	if err != nil {
		t.Fatal(err)
	}
	if mem.Fd() < 0 {
		t.Fatalf("invalid fd")
	}
	if mem.Raw() == nil {
		t.Fatal("nil raw pointer")
	}
	if mem.Size() != size {
		t.Fatal("wrong size")
	}
	if mem.Offset() != 0 {
		t.Fatal("wrong offset")
	}
}

func TestDestroyPosixSharedMemory(t *testing.T) {
	mem, err := NewPosix(size, flag, perm)
	if err != nil {
		t.Fatal(err)
	}
	err = mem.Close()
	if err != nil {
		t.Fatal(err)
	}
}

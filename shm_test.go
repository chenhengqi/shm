package shm

import (
	"os"
	"testing"
)

const (
	size = 1024 * 1024
	flag = os.O_RDWR | os.O_CREATE | os.O_EXCL
	perm = 0600
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

func TestRemovePosixSharedMemory(t *testing.T) {
	mem, err := NewPosix(size, flag, perm)
	if err != nil {
		t.Fatal(err)
	}
	err = mem.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateSystemVSharedMemory(t *testing.T) {
	mem, err := NewSystemVPrivate(size, IPCCreate|IPCExclusive, perm)
	if err != nil {
		t.Fatal(err)
	}
	if mem.Fd() != ^uintptr(0) {
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

func TestRemoveSystemVSharedMemory(t *testing.T) {
	mem, err := NewSystemVPrivate(size, IPCCreate, perm)
	if err != nil {
		t.Fatal(err)
	}
	err = mem.Close()
	if err != nil {
		t.Fatal(err)
	}
}

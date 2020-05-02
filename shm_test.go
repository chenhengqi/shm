package shm

import (
	"io"
	"os"
	"testing"
)

const (
	size = 1024 * 1024
	flag = os.O_RDWR | os.O_CREATE | os.O_EXCL
	perm = 0666
)

func TestPosixSharedMemory(t *testing.T) {
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

	data := []byte("hello world")
	n, err := mem.Write(data)
	if n != len(data) {
		t.Fatal("short write")
	}
	if mem.Offset() != int64(len(data)) {
		t.Fatal("invalid offset")
	}

	offset, err := mem.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatal(err)
	}
	if offset != 0 {
		t.Fatal("seek failed")
	}

	buf := make([]byte, len(data))
	n, err = mem.Read(buf)
	if n != len(buf) {
		t.Fatal("read failed")
	}
	if string(buf) != string(data) {
		t.Fatal("read failed")
	}

	err = mem.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSystemVSharedMemory(t *testing.T) {
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

	data := []byte("hello world")
	n, err := mem.Write(data)
	if n != len(data) {
		t.Fatal("short write")
	}
	if mem.Offset() != int64(len(data)) {
		t.Fatal("invalid offset")
	}

	offset, err := mem.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatal(err)
	}
	if offset != 0 {
		t.Fatal("seek failed")
	}

	buf := make([]byte, len(data))
	n, err = mem.Read(buf)
	if n != len(buf) {
		t.Fatal("read failed")
	}
	if string(buf) != string(data) {
		t.Fatal("read failed")
	}

	err = mem.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRandomName(t *testing.T) {
	name1 := randomName()
	name2 := randomName()
	if len(name1) != 26 {
		t.Fatal("unexpected name length")
	}
	if name1 == name2 {
		t.Fatal("same name")
	}
}

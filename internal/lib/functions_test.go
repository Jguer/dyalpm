package lib

import (
	"testing"
	"unsafe"
)

func TestPtrToStringCopies(t *testing.T) {
	buf := []byte("hello\x00")
	s := PtrToString(uintptr(unsafe.Pointer(&buf[0])))
	buf[0] = 'j'
	if s != "hello" {
		t.Fatalf("PtrToString aliases C memory: got %q after mutation, want %q", s, "hello")
	}
}

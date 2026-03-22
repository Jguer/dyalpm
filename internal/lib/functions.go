package lib

import (
	"unsafe"
)

// CString converts a Go string to a null-terminated byte slice.
// The caller must keep the slice alive during the C call and use &buf[0] as the pointer.
func CString(s string) []byte {
	if s == "" {
		return []byte{0}
	}

	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = 0
	return buf
}

// PtrToString converts a C string pointer to a Go string
// This function finds the null terminator to determine the length
func PtrToString(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}
	base := unsafe.Pointer(ptr)
	start := (*byte)(base)
	if start == nil {
		return ""
	}
	// Find null terminator
	n := 0
	for *(*byte)(unsafe.Add(base, n)) != 0 {
		n++
		if n > 1024*1024 { // Safety limit
			return ""
		}
	}
	if n == 0 {
		return ""
	}
	return unsafe.String(start, n)
}

// PtrToStringWithLen converts a C string pointer with known length to Go string
func PtrToStringWithLen(ptr uintptr, length int) string {
	if ptr == 0 || length == 0 {
		return ""
	}
	base := unsafe.Pointer(ptr)
	return unsafe.String((*byte)(base), length)
}

// BoolToInt converts a Go bool to C int (0 or 1)
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// IntToBool converts a C int to Go bool
func IntToBool(i int32) bool {
	return i != 0
}

package alpmList

import (
	"unsafe"

	"github.com/Jguer/dyalpm/internal/lib"
)

// List represents an ALPM list node
type List struct {
	ptr uintptr
}

// NewList creates a new List wrapper from a C pointer
func NewList(ptr uintptr) *List {
	if ptr == 0 {
		return nil
	}
	return &List{ptr: ptr}
}

// Ptr returns the underlying C pointer
func (l *List) Ptr() uintptr {
	if l == nil {
		return 0
	}
	return l.ptr
}

// Data returns the data pointer from the list node
func (l *List) Data() uintptr {
	if l.ptr == 0 {
		return 0
	}
	base := unsafe.Pointer(l.ptr)
	return *(*uintptr)(base)
}

// Next returns the next node in the list
func (l *List) Next() *List {
	if l.ptr == 0 {
		return nil
	}
	// alpm_list_t structure: data (0), prev (8), next (16)
	// We want offset 16 (2*ptrSize) for next
	base := unsafe.Pointer(l.ptr)
	nextPtr := *(*uintptr)(unsafe.Add(base, 2*unsafe.Sizeof(uintptr(0))))
	if nextPtr == 0 {
		return nil
	}
	return &List{ptr: nextPtr}
}

// Free frees the list structure (but not the data)
func (l *List) Free() {
	if l == nil || l.ptr == 0 {
		return
	}
	if lib.AlpmListFree != nil {
		lib.AlpmListFree(l.ptr)
	}
}

// Add adds data to the list and returns the new list head
func Add(l *List, data uintptr) *List {
	if lib.AlpmListAdd == nil {
		return nil
	}
	var ptr uintptr
	if l != nil {
		ptr = l.ptr
	}
	r1 := lib.AlpmListAdd(ptr, data)
	return NewList(r1)
}

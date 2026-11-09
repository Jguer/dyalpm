package alpmList

import (
	"testing"
	"unsafe"
)

// alpmListNode represents the C struct alpm_list_t layout for testing
// struct alpm_list_t {
//
//	void *data;
//	alpm_list_t *prev;
//	alpm_list_t *next;
//
// }
type alpmListNode struct {
	data uintptr
	prev uintptr
	next uintptr
}

// createTestList creates a linked list of nodes with the given data values.
// Returns the head of the list.
func createTestList(data []uintptr) (*List, []alpmListNode) {
	if len(data) == 0 {
		return nil, nil
	}

	nodes := make([]alpmListNode, len(data))
	for i := range nodes {
		nodes[i].data = data[i]
		if i > 0 {
			nodes[i].prev = uintptr(unsafe.Pointer(&nodes[i-1]))
		}
		if i < len(data)-1 {
			nodes[i].next = uintptr(unsafe.Pointer(&nodes[i+1]))
		}
	}

	return NewList(uintptr(unsafe.Pointer(&nodes[0]))), nodes
}

func TestNewList_Nil(t *testing.T) {
	l := NewList(0)
	if l != nil {
		t.Error("NewList(0) should return nil")
	}
}

func TestNewList_Valid(t *testing.T) {
	node := alpmListNode{data: 42}
	l := NewList(uintptr(unsafe.Pointer(&node)))
	if l == nil {
		t.Fatal("NewList should not return nil for valid pointer")
	}
	if l.Ptr() == 0 {
		t.Error("Ptr() should not be 0")
	}
}

func TestList_Ptr_Nil(t *testing.T) {
	var l *List = nil
	if l.Ptr() != 0 {
		t.Error("nil List.Ptr() should return 0")
	}
}

func TestList_Data(t *testing.T) {
	node := alpmListNode{data: 0xDEADBEEF}
	l := NewList(uintptr(unsafe.Pointer(&node)))

	data := l.Data()
	if data != 0xDEADBEEF {
		t.Errorf("Data() = %#x, want %#x", data, 0xDEADBEEF)
	}
}

func TestList_Next(t *testing.T) {
	l, _ := createTestList([]uintptr{100, 200, 300})

	next := l.Next()
	if next == nil {
		t.Fatal("Next() should not be nil")
	}
	if next.Data() != 200 {
		t.Errorf("Next().Data() = %d, want 200", next.Data())
	}

	next2 := next.Next()
	if next2 == nil {
		t.Fatal("Second Next() should not be nil")
	}
	if next2.Data() != 300 {
		t.Errorf("Next().Next().Data() = %d, want 300", next2.Data())
	}

	next3 := next2.Next()
	if next3 != nil {
		t.Error("Third Next() should be nil (end of list)")
	}
}

func TestList_SingleElement(t *testing.T) {
	l, _ := createTestList([]uintptr{42})

	if l.Data() != 42 {
		t.Errorf("Data() = %d, want 42", l.Data())
	}
	if l.Next() != nil {
		t.Error("single element Next() should be nil")
	}
}

func TestList_ZeroPtr_Data(t *testing.T) {
	l := &List{ptr: 0}
	data := l.Data()
	if data != 0 {
		t.Errorf("zero ptr Data() = %d, want 0", data)
	}
}

func TestList_ZeroPtr_Next(t *testing.T) {
	l := &List{ptr: 0}
	next := l.Next()
	if next != nil {
		t.Error("zero ptr Next() should be nil")
	}
}

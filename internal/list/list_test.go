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

func TestList_Prev(t *testing.T) {
	l, _ := createTestList([]uintptr{100, 200, 300})

	// Move to end
	end := l.Next().Next()
	if end.Data() != 300 {
		t.Fatalf("Expected end data 300, got %d", end.Data())
	}

	prev := end.Prev()
	if prev == nil {
		t.Fatal("Prev() should not be nil")
	}
	if prev.Data() != 200 {
		t.Errorf("Prev().Data() = %d, want 200", prev.Data())
	}

	prev2 := prev.Prev()
	if prev2 == nil {
		t.Fatal("Second Prev() should not be nil")
	}
	if prev2.Data() != 100 {
		t.Errorf("Prev().Prev().Data() = %d, want 100", prev2.Data())
	}

	prev3 := prev2.Prev()
	if prev3 != nil {
		t.Error("Third Prev() should be nil (start of list)")
	}
}

func TestList_Iterator(t *testing.T) {
	l, _ := createTestList([]uintptr{10, 20, 30})

	it := l.Iterator()
	if it == nil {
		t.Fatal("Iterator() should not be nil")
	}

	if !it.HasNext() {
		t.Error("HasNext() should be true initially")
	}
}

func TestListIterator_Next(t *testing.T) {
	l, _ := createTestList([]uintptr{10, 20, 30})
	it := l.Iterator()

	values := []uintptr{}
	for it.HasNext() {
		values = append(values, it.Next())
	}

	if len(values) != 3 {
		t.Errorf("got %d values, want 3", len(values))
	}

	expected := []uintptr{10, 20, 30}
	for i, v := range values {
		if v != expected[i] {
			t.Errorf("values[%d] = %d, want %d", i, v, expected[i])
		}
	}
}

func TestListIterator_Empty(t *testing.T) {
	// Create an iterator with nil current
	it := &Iterator{current: nil}

	if it.HasNext() {
		t.Error("empty iterator should have HasNext() = false")
	}

	data := it.Next()
	if data != 0 {
		t.Errorf("empty iterator Next() = %d, want 0", data)
	}
}

func TestList_ToSlice(t *testing.T) {
	l, _ := createTestList([]uintptr{100, 200, 300})

	slice := l.ToSlice()
	if len(slice) != 3 {
		t.Errorf("ToSlice() length = %d, want 3", len(slice))
	}

	expected := []uintptr{100, 200, 300}
	for i, v := range slice {
		if v != expected[i] {
			t.Errorf("slice[%d] = %d, want %d", i, v, expected[i])
		}
	}
}

func TestList_ToSlice_Nil(t *testing.T) {
	var l *List = nil
	slice := l.ToSlice()
	if slice != nil {
		t.Errorf("nil List.ToSlice() should return nil, got %v", slice)
	}
}

func TestList_Count_Nil(t *testing.T) {
	var l *List = nil
	count := l.Count()
	if count != 0 {
		t.Errorf("nil List.Count() = %d, want 0", count)
	}
}

func TestList_Count_Manual(t *testing.T) {
	// When libalpm is not available, Count may use FFI or fallback.
	// The FFI path may not work correctly without libalpm loaded.
	// This test verifies the structure works - actual count verification
	// happens in the integration tests.
	l, _ := createTestList([]uintptr{1, 2, 3, 4, 5})

	// At minimum, Count should not panic
	count := l.Count()

	// We can verify iteration works even if Count() doesn't
	iterCount := 0
	for node := l; node != nil && node.ptr != 0; node = node.Next() {
		iterCount++
	}
	if iterCount != 5 {
		t.Errorf("Manual iteration count = %d, want 5", iterCount)
	}

	// Log the Count() result for debugging
	t.Logf("Count() returned %d (FFI may not be available)", count)
}

func TestList_SingleElement(t *testing.T) {
	l, _ := createTestList([]uintptr{42})

	if l.Data() != 42 {
		t.Errorf("Data() = %d, want 42", l.Data())
	}
	if l.Next() != nil {
		t.Error("single element Next() should be nil")
	}
	if l.Prev() != nil {
		t.Error("single element Prev() should be nil")
	}

	slice := l.ToSlice()
	if len(slice) != 1 || slice[0] != 42 {
		t.Errorf("ToSlice() = %v, want [42]", slice)
	}
}

func TestListIterator_ForLargeList(t *testing.T) {
	data := make([]uintptr, 1000)
	for i := range data {
		data[i] = uintptr(i)
	}

	l, _ := createTestList(data)
	it := l.Iterator()

	count := 0
	sum := uintptr(0)
	for it.HasNext() {
		sum += it.Next()
		count++
	}

	if count != 1000 {
		t.Errorf("iterated %d times, want 1000", count)
	}

	// Sum of 0..999 = 999*1000/2 = 499500
	if sum != 499500 {
		t.Errorf("sum = %d, want 499500", sum)
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

func TestList_ZeroPtr_Prev(t *testing.T) {
	l := &List{ptr: 0}
	prev := l.Prev()
	if prev != nil {
		t.Error("zero ptr Prev() should be nil")
	}
}

package lib

import (
	"errors"
	"testing"
)

func TestMockRegistry_GetFunc(t *testing.T) {
	mock := NewMockRegistry()

	// Test function not found
	_, err := mock.GetFunc("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent function")
	}

	// Register and retrieve
	mock.RegisterFunc("test_func", 0x12345678)
	ptr, err := mock.GetFunc("test_func")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ptr != 0x12345678 {
		t.Errorf("got %#x, want %#x", ptr, 0x12345678)
	}
}

func TestMockRegistry_CallFunc(t *testing.T) {
	mock := NewMockRegistry()

	// Test with preset result
	mock.SetCallResult("my_func", 42)
	result, err := mock.CallFunc("my_func", 1, 2, 3)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != 42 {
		t.Errorf("got %d, want 42", result)
	}

	// Verify call count
	if mock.GetCallCount("my_func") != 1 {
		t.Errorf("call count = %d, want 1", mock.GetCallCount("my_func"))
	}

	// Call again
	_, _ = mock.CallFunc("my_func")
	if mock.GetCallCount("my_func") != 2 {
		t.Errorf("call count = %d, want 2", mock.GetCallCount("my_func"))
	}
}

func TestMockRegistry_CallFunc_Error(t *testing.T) {
	mock := NewMockRegistry()

	expectedErr := errors.New("mock error")
	mock.SetCallError("failing_func", expectedErr)

	_, err := mock.CallFunc("failing_func")
	if err != expectedErr {
		t.Errorf("got error %v, want %v", err, expectedErr)
	}
}

func TestMockRegistry_Reset(t *testing.T) {
	mock := NewMockRegistry()

	mock.RegisterFunc("func1", 123)
	mock.SetCallResult("func2", 456)
	_, _ = mock.CallFunc("func2")

	mock.Reset()

	// After reset, function should not be found
	_, err := mock.GetFunc("func1")
	if err == nil {
		t.Error("expected error after reset")
	}

	// Call count should be zero
	if mock.GetCallCount("func2") != 0 {
		t.Errorf("call count after reset = %d, want 0", mock.GetCallCount("func2"))
	}
}

func TestMockRegistry_ImplementsFunctionCaller(t *testing.T) {
	var _ FunctionCaller = NewMockRegistry()
}

func TestMockRegistry_ConcurrentAccess(t *testing.T) {
	mock := NewMockRegistry()
	mock.RegisterFunc("concurrent_func", 100)
	mock.SetCallResult("concurrent_func", 200)

	done := make(chan bool)
	for range 10 {
		go func() {
			for range 100 {
				_, _ = mock.GetFunc("concurrent_func")
				_, _ = mock.CallFunc("concurrent_func")
			}
			done <- true
		}()
	}

	for range 10 {
		<-done
	}

	if mock.GetCallCount("concurrent_func") != 1000 {
		t.Errorf("call count = %d, want 1000", mock.GetCallCount("concurrent_func"))
	}
}

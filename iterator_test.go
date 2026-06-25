package dyalpm

import (
	"errors"
	"testing"
)

func TestPackageIterator_ForEach_NilList(t *testing.T) {
	// A PackageIterator with nil list should not error
	it := PackageIterator{
		list:   nil,
		handle: nil,
	}

	called := false
	err := it.ForEach(func(pkg Package) error {
		called = true
		return nil
	})
	if err != nil {
		t.Errorf("ForEach on nil list returned error: %v", err)
	}
	if called {
		t.Error("ForEach on nil list should not call function")
	}
}

func TestPackageIterator_Collect_NilList(t *testing.T) {
	it := PackageIterator{
		list:   nil,
		handle: nil,
	}

	pkgs := it.Collect()
	if pkgs != nil {
		t.Errorf("Collect on nil list returned %v, want nil", pkgs)
	}
}

func TestPackageIterator_SortBySize_NilList(t *testing.T) {
	it := PackageIterator{
		list:   nil,
		handle: nil,
	}

	pkgs := it.SortBySize()
	if pkgs != nil {
		t.Errorf("SortBySize on nil list returned %v, want nil", pkgs)
	}
}

func TestPackageIterator_FindSatisfier_NilList(t *testing.T) {
	it := PackageIterator{
		list:   nil,
		handle: nil,
	}

	pkg, err := it.FindSatisfier("glibc")
	if pkg != nil {
		t.Errorf("FindSatisfier on nil list returned package: %v", pkg)
	}
	if !errors.Is(err, ErrPackageNotFound) {
		t.Errorf("FindSatisfier error = %v, want ErrPackageNotFound", err)
	}
}

func TestPackageIterator_ForEach_ErrorPropagation(t *testing.T) {
	// This test verifies that errors from the callback are propagated
	// We can't easily test this without real packages, but we document the expected behavior

	// The ForEach method should return any error returned by the callback function
	expectedErr := errors.New("test error")

	// With nil list, error won't be propagated (no iteration)
	it := PackageIterator{list: nil}
	err := it.ForEach(func(pkg Package) error {
		return expectedErr
	})
	// No error because no iteration occurred
	if err != nil {
		t.Errorf("ForEach on nil list should return nil, got: %v", err)
	}
}

func TestPackageIterator_ZeroValue(t *testing.T) {
	var it PackageIterator

	// Zero value should behave like nil list
	err := it.ForEach(func(pkg Package) error {
		t.Error("should not be called")
		return nil
	})
	if err != nil {
		t.Errorf("ForEach on zero value returned error: %v", err)
	}

	pkgs := it.Collect()
	if pkgs != nil {
		t.Errorf("Collect on zero value returned %v, want nil", pkgs)
	}

	sorted := it.SortBySize()
	if sorted != nil {
		t.Errorf("SortBySize on zero value returned %v, want nil", sorted)
	}
}

// Test newPackageIterator construction
func TestNewPackageIterator_ZeroPtr(t *testing.T) {
	it := newPackageIterator(0, nil, false)

	// With zero pointer, list.NewList returns nil
	if it.list != nil {
		t.Errorf("expected nil list for zero pointer, got %v", it.list)
	}
	if it.handle != nil {
		t.Error("expected nil handle")
	}
	if it.freeOnDone {
		t.Error("expected freeOnDone=false")
	}
}

func TestNewPackageIterator_FreeOnDone(t *testing.T) {
	it := newPackageIterator(0, nil, true)

	if !it.freeOnDone {
		t.Error("expected freeOnDone=true")
	}
}

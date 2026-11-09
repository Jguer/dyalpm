package dyalpm

import (
	"cmp"
	"slices"

	alpmlist "github.com/Jguer/dyalpm/internal/list"
)

// PackageIterator provides lazy iteration over a libalpm-owned package list,
// or over an already collected slice. It can be iterated any number of times.
type PackageIterator struct {
	list   *alpmlist.List
	pkgs   []Package
	handle *handle
}

// newPackageIterator wraps a list owned by libalpm; the iterator never frees it.
func newPackageIterator(listPtr uintptr, h *handle) PackageIterator {
	return PackageIterator{
		list:   alpmlist.NewList(listPtr),
		handle: h,
	}
}

// ForEach iterates over the packages.
func (it PackageIterator) ForEach(fn func(Package) error) error {
	for _, pkg := range it.pkgs {
		if err := fn(pkg); err != nil {
			return err
		}
	}
	for item := it.list; item != nil && item.Ptr() != 0; item = item.Next() {
		pkgPtr := item.Data()
		if pkgPtr == 0 {
			continue
		}
		if err := fn(newPackage(pkgPtr, it.handle)); err != nil {
			return err
		}
	}
	return nil
}

// Collect returns all packages in the iterator as a slice.
func (it PackageIterator) Collect() []Package {
	var pkgs []Package
	_ = it.ForEach(func(pkg Package) error {
		pkgs = append(pkgs, pkg)
		return nil
	})
	return pkgs
}

// FindSatisfier finds the first package satisfying a depstring.
func (it PackageIterator) FindSatisfier(depstring string) (Package, error) {
	pkg := FindSatisfier(it.Collect(), depstring)
	if pkg == nil {
		return nil, ErrPackageNotFound
	}
	return pkg, nil
}

// SortBySize returns packages sorted by install size (descending).
func (it PackageIterator) SortBySize() []Package {
	pkgs := it.Collect()
	slices.SortFunc(pkgs, func(a, b Package) int {
		return cmp.Compare(b.ISize(), a.ISize())
	})
	return pkgs
}

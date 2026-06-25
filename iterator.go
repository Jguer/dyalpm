package dyalpm

import (
	"cmp"
	"slices"

	alpmlist "github.com/Jguer/dyalpm/internal/list"
)

// PackageIterator provides lazy iteration over an ALPM package list.
type PackageIterator struct {
	list       *alpmlist.List
	handle     *handle
	freeOnDone bool
}

func newPackageIterator(listPtr uintptr, h *handle, freeOnDone bool) PackageIterator {
	return PackageIterator{
		list:       alpmlist.NewList(listPtr),
		handle:     h,
		freeOnDone: freeOnDone,
	}
}

// ForEach iterates over packages lazily. If the underlying list must be freed,
// it will be freed after iteration.
func (it PackageIterator) ForEach(fn func(Package) error) error {
	if it.list == nil {
		return nil
	}
	if it.freeOnDone {
		defer it.list.Free()
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

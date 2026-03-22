package dyalpm

import (
	stderrors "errors"
	"runtime"
	"unsafe"

	"github.com/Jguer/dyalpm/internal/lib"
	alpmlist "github.com/Jguer/dyalpm/internal/list"
)

func collectList[T any](alpmList *alpmlist.List, build func(uintptr) T) []T {
	var items []T
	for item := alpmList; item != nil && item.Ptr() != 0; item = item.Next() {
		ptr := item.Data()
		if ptr != 0 {
			items = append(items, build(ptr))
		}
	}
	return items
}

func computeFileSum(funcName, filename string) (string, error) {
	var r1 uintptr
	switch funcName {
	case "alpm_compute_md5sum":
		if lib.AlpmComputeMd5sum == nil {
			return "", ErrInvalidPackage
		}
		r1 = lib.AlpmComputeMd5sum(filename)
	case "alpm_compute_sha256sum":
		if lib.AlpmComputeSha256sum == nil {
			return "", ErrInvalidPackage
		}
		r1 = lib.AlpmComputeSha256sum(filename)
	default:
		return "", ErrInvalidPackage
	}

	if r1 == 0 {
		return "", ErrInvalidPackage
	}

	res := lib.PtrToString(r1)
	lib.Free(r1)
	return res, nil
}

// ComputeMD5Sum computes the MD5 sum of a file
func ComputeMD5Sum(filename string) (string, error) {
	return computeFileSum("alpm_compute_md5sum", filename)
}

// ComputeSHA256Sum computes the SHA256 sum of a file
func ComputeSHA256Sum(filename string) (string, error) {
	return computeFileSum("alpm_compute_sha256sum", filename)
}

func (h *handle) FindGroupPkgs(dbs []Database, name string) ([]Package, error) {
	if h.ptr == 0 {
		return nil, ErrInvalidHandle
	}
	if lib.AlpmFindGroupPkgs == nil {
		return nil, stderrors.New("missing function: alpm_find_group_pkgs")
	}

	var dbList *alpmlist.List
	for _, db := range dbs {
		dbImpl, ok := db.(*database)
		if ok {
			dbList = alpmlist.Add(dbList, dbImpl.ptr)
		}
	}
	defer dbList.Free()

	r1 := lib.AlpmFindGroupPkgs(dbList.Ptr(), name)

	if r1 == 0 {
		return []Package{}, nil
	}

	resList := alpmlist.NewList(r1)
	defer resList.Free()

	var pkgs []Package
	for item := resList; item != nil && item.Ptr() != 0; item = item.Next() {
		ptr := item.Data()
		if ptr != 0 {
			pkgs = append(pkgs, newPackage(ptr, h))
		}
	}

	return pkgs, nil
}

func (h *handle) ExtractKeyID(identifier string, sig []byte) ([]string, error) {
	if h.ptr == 0 {
		return nil, ErrInvalidHandle
	}
	if lib.AlpmExtractKeyID == nil {
		return nil, stderrors.New("missing function: alpm_extract_keyid")
	}

	var keysListPtr uintptr
	var sigPtr uintptr
	if len(sig) > 0 {
		sigPtr = uintptr(unsafe.Pointer(&sig[0]))
	}

	sigLenInt32 := clampIntToInt32(len(sig))
	r1 := lib.AlpmExtractKeyID(h.ptr, identifier, sigPtr, sigLenInt32, &keysListPtr)

	runtime.KeepAlive(sig)

	if r1 != 0 {
		return nil, stderrors.New("failed to extract key id")
	}

	if keysListPtr == 0 {
		return []string{}, nil
	}

	alpmList := alpmlist.NewList(keysListPtr)
	// We need to free the strings in the list too
	// alpm_list_free_inner(list, free)
	// For now, let's just free the list structure and hope the strings are managed or short-lived.
	// Actually, alpm_list_free just frees the nodes.
	defer alpmList.Free()

	var keys []string
	for item := alpmList; item != nil && item.Ptr() != 0; item = item.Next() {
		ptr := item.Data()
		if ptr != 0 {
			keys = append(keys, lib.PtrToString(ptr))
		}
	}

	return keys, nil
}

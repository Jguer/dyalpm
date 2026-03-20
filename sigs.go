package dyalpm

import (
	stderrors "errors"
	"unsafe"

	"github.com/Jguer/dyalpm/internal/lib"
)

// SigStatus represents the status of a signature
type SigStatus int

const (
	SigStatusValid SigStatus = iota
	SigStatusInvalid
	SigStatusSigExpired
	SigStatusKeyExpired
	SigStatusKeyUnknown
	SigStatusKeyDisabled
	SigStatusKeyRevoked
)

// SigValidity represents the validity of a signature
type SigValidity int

const (
	SigValidityFull SigValidity = iota
	SigValidityMarginal
	SigValidityNever
	SigValidityUnknown
)

// SigResult represents a single signature result
type SigResult struct {
	KeyID    string
	Status   SigStatus
	Validity SigValidity
}

// SigList represents a list of signature results
type SigList struct {
	Results []SigResult
}

// Internal structure for alpm_siglist_t
type alpmSiglistT struct {
	Count   uintptr
	Results uintptr
}

// Internal structure for alpm_sigresult_t
type alpmSigresultT struct {
	KeyID    uintptr
	Status   int
	Validity int
}

func decodeSigList(ptr uintptr) SigList {
	if ptr == 0 {
		return SigList{}
	}

	base := unsafe.Pointer(ptr)
	sigList := (*alpmSiglistT)(base)
	if sigList.Count == 0 || sigList.Results == 0 {
		return SigList{}
	}

	var results []SigResult
	resultSize := unsafe.Sizeof(alpmSigresultT{})
	resultsBase := unsafe.Pointer(sigList.Results)

	for i := uintptr(0); i < sigList.Count; i++ {
		res := (*alpmSigresultT)(unsafe.Add(resultsBase, i*resultSize))

		results = append(results, SigResult{
			KeyID:    lib.PtrToString(res.KeyID),
			Status:   SigStatus(res.Status),
			Validity: SigValidity(res.Validity),
		})
	}

	return SigList{Results: results}
}

func checkPGPSignature(ptr uintptr, handle *handle, funcName string) (SigList, error) {
	var fn func(uintptr, uintptr) int32
	switch funcName {
	case "alpm_db_check_pgp_signature":
		fn = lib.AlpmDBCheckPGPSignature
	case "alpm_pkg_check_pgp_signature":
		fn = lib.AlpmPkgCheckPGPSignature
	default:
		return SigList{}, stderrors.New("missing function: " + funcName)
	}

	if fn == nil {
		return SigList{}, stderrors.New("missing function: " + funcName)
	}

	var sigList alpmSiglistT
	sigListPtr := uintptr(unsafe.Pointer(&sigList))

	result := fn(ptr, sigListPtr)
	decoded := decodeSigList(sigListPtr)

	cleanupErr := handle.SigListCleanup(sigListPtr)

	if result != 0 {
		err := stderrors.New("signature check failed")
		if cleanupErr != nil {
			return decoded, stderrors.Join(err, cleanupErr)
		}
		return decoded, err
	}
	if cleanupErr != nil {
		return decoded, cleanupErr
	}

	return decoded, nil
}

// SigListCleanup cleans up an ALPM signature list
func (h *handle) SigListCleanup(siglistPtr uintptr) error {
	if h.ptr == 0 {
		return ErrInvalidHandle
	}
	if lib.AlpmSiglistCleanup == nil {
		return stderrors.New("missing function: alpm_siglist_cleanup")
	}

	if lib.AlpmSiglistCleanup(siglistPtr, h.ptr) != 0 {
		return stderrors.New("failed to clean up signature list")
	}
	return nil
}

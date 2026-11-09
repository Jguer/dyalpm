package dyalpm

import (
	stderrors "errors"
	"unsafe"

	"github.com/Jguer/dyalpm/internal/lib"
)

// SigStatus represents the status of a signature (alpm_sigstatus_t).
type SigStatus int

const (
	SigStatusValid SigStatus = iota
	SigStatusKeyExpired
	SigStatusSigExpired
	SigStatusKeyUnknown
	SigStatusKeyDisabled
	SigStatusInvalid
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
	KeyID    string // key fingerprint
	Status   SigStatus
	Validity SigValidity
}

// SigList represents a list of signature results
type SigList struct {
	Results []SigResult
}

// alpmSiglistT mirrors alpm_siglist_t.
type alpmSiglistT struct {
	Count   uintptr
	Results uintptr
}

// alpmSigresultT mirrors alpm_sigresult_t, which embeds a whole alpm_pgpkey_t.
type alpmSigresultT struct {
	KeyData     uintptr
	Fingerprint uintptr
	UID         uintptr
	Name        uintptr
	Email       uintptr
	Created     int64
	Expires     int64
	Length      uint32
	Revoked     uint32
	Status      int32
	Validity    int32
}

func decodeSigList(sigList *alpmSiglistT) SigList {
	if sigList.Count == 0 || sigList.Results == 0 {
		return SigList{}
	}

	results := unsafe.Slice((*alpmSigresultT)(unsafe.Pointer(sigList.Results)), sigList.Count)
	out := make([]SigResult, len(results))
	for i, res := range results {
		out[i] = SigResult{
			KeyID:    lib.PtrToString(res.Fingerprint),
			Status:   SigStatus(res.Status),
			Validity: SigValidity(res.Validity),
		}
	}

	return SigList{Results: out}
}

func checkPGPSignature(ptr uintptr, fn func(uintptr, unsafe.Pointer) int32, funcName string) (SigList, error) {
	if fn == nil {
		return SigList{}, stderrors.New("missing function: " + funcName)
	}
	if lib.AlpmSiglistCleanup == nil {
		return SigList{}, stderrors.New("missing function: alpm_siglist_cleanup")
	}

	var sigList alpmSiglistT
	result := fn(ptr, unsafe.Pointer(&sigList))
	decoded := decodeSigList(&sigList)

	var err error
	if result != 0 {
		err = stderrors.New("signature check failed")
	}
	if lib.AlpmSiglistCleanup(unsafe.Pointer(&sigList)) != 0 {
		err = stderrors.Join(err, stderrors.New("failed to clean up signature list"))
	}

	return decoded, err
}

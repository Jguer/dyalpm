package dyalpm

import (
	"runtime"
	"testing"
	"unsafe"
)

func TestDecodeSigList_Layout(t *testing.T) {
	if size := unsafe.Sizeof(alpmSigresultT{}); size != 72 {
		t.Fatalf("alpm_sigresult_t is 72 bytes on 64-bit, mirror is %d", size)
	}

	fpr := []byte("ABCD\x00")
	results := []alpmSigresultT{
		{Status: 0, Validity: 0},
		{Fingerprint: uintptr(unsafe.Pointer(&fpr[0])), Status: 1, Validity: 3},
	}
	list := alpmSiglistT{Count: 2, Results: uintptr(unsafe.Pointer(&results[0]))}

	got := decodeSigList(&list).Results
	runtime.KeepAlive(results)
	runtime.KeepAlive(fpr)
	if len(got) != 2 {
		t.Fatalf("got %d results, want 2", len(got))
	}
	if got[1].KeyID != "ABCD" || got[1].Status != SigStatusKeyExpired || got[1].Validity != SigValidityUnknown {
		t.Fatalf("second result decoded wrong: %+v", got[1])
	}
}

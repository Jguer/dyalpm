package dyalpm

import (
	stderrors "errors"

	"github.com/Jguer/dyalpm/internal/lib"
)

// Version returns the library version string
func Version() string {
	if err := lib.EnsureAlpmLoaded(); err != nil {
		return ""
	}
	if lib.AlpmVersion == nil {
		return ""
	}
	return lib.PtrToString(lib.AlpmVersion())
}

// Capabilities returns the library capabilities
func Capabilities() (CapabilitiesMask, error) {
	if err := lib.EnsureAlpmLoaded(); err != nil {
		return 0, err
	}
	if lib.AlpmCapabilities == nil {
		return 0, stderrors.New("missing function: alpm_capabilities")
	}
	return CapabilitiesMask(lib.AlpmCapabilities()), nil
}

// CapabilitiesMask represents the library capabilities bitmask
type CapabilitiesMask uint64

const (
	CapNLS        CapabilitiesMask = 1 << 0
	CapDownloader CapabilitiesMask = 1 << 1
	CapSignatures CapabilitiesMask = 1 << 2
)

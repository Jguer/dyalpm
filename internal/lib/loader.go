package lib

import (
	"errors"
	"sync"

	"github.com/ebitengine/purego"
)

const (
	libalpmPath         = "libalpm.so.16"
	libalpmPathFallback = "libalpm.so"
)

var (
	loadMu     sync.Mutex
	libalpm    uintptr
	libLoaded  bool
	libalpmErr error
)

// LoadLibrary loads the ALPM library. It's safe to call multiple times.
func LoadLibrary(libalpmPath string) error {
	loadMu.Lock()
	defer loadMu.Unlock()
	return loadALPMLibraryUnlocked(libalpmPath)
}

func loadALPMLibraryUnlocked(primaryPath string) error {
	if libLoaded {
		return libalpmErr
	}

	handle, err := purego.Dlopen(primaryPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		handle, err = purego.Dlopen(libalpmPathFallback, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	}
	if err != nil {
		libalpmErr = err
		return err
	}

	libalpm = handle
	libalpmErr = nil
	libLoaded = true
	registerAlpmFuncs(libalpm)
	registerLibcFuncs()

	return nil
}

// LoadAlpm loads libalpm (primary soname, then libalpm.so). Retries after a failed load.
func LoadAlpm(libalpmPath string) error {
	loadMu.Lock()
	defer loadMu.Unlock()
	if libLoaded {
		return libalpmErr
	}
	return loadALPMLibraryUnlocked(libalpmPath)
}

func EnsureAlpmLoaded() error {
	_, err := getLibrary()
	return err
}

// getLibrary returns the loaded library handle. It will attempt to load if not already loaded.
func getLibrary() (uintptr, error) {
	if err := LoadAlpm(libalpmPath); err != nil {
		return 0, err
	}
	if libalpm == 0 {
		return 0, errors.New("library not loaded")
	}
	return libalpm, nil
}

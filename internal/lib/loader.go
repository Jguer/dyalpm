package lib

import (
	"sync"

	"github.com/ebitengine/purego"
)

const (
	libalpmPath         = "libalpm.so.16"
	libalpmPathFallback = "libalpm.so"
)

var (
	loadMu  sync.Mutex
	libalpm uintptr
)

// EnsureAlpmLoaded loads libalpm (primary soname, then libalpm.so) once.
// A failed load is retried on the next call.
func EnsureAlpmLoaded() error {
	loadMu.Lock()
	defer loadMu.Unlock()
	if libalpm != 0 {
		return nil
	}

	handle, err := purego.Dlopen(libalpmPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		handle, err = purego.Dlopen(libalpmPathFallback, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	}
	if err != nil {
		return err
	}

	registerAlpmFuncs(handle)
	registerLibcFuncs()
	libalpm = handle
	return nil
}
